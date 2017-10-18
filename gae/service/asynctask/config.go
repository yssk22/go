package asynctask

import (
	"fmt"

	"net/url"

	"time"

	"context"

	"github.com/speedland/go/gae/taskqueue"
	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/lazy"
	"github.com/speedland/go/x/xlog"
	"github.com/speedland/go/x/xruntime"
	"github.com/speedland/go/x/xtime"
)

// Config is an configuration object to define AsyncTask endpoints
type Config struct {
	key         string
	queue       *taskqueue.PushQueue
	logic       Logic
	schedule    string
	description string
}

// NewConfig returns a new *Config named by `key`` to execute the async logic on top of the queue.
func NewConfig(key string, queue *taskqueue.PushQueue, logic Logic) *Config {
	caller := xruntime.CaptureCaller()
	desc := fmt.Sprintf("defined in %s:%d", caller.FullFilePath, caller.LineNumber)
	if key == "" {
		panic(fmt.Errorf("asynctask.Config key must not be empty"))
	}
	return &Config{
		key:         key,
		queue:       queue,
		logic:       logic,
		description: desc,
	}
}

// GetKey returns a key string for the async task
func (c *Config) GetKey() string {
	return c.key
}

// WithDescription sets the description for the async task
func (c *Config) WithDescription(description string) *Config {
	c.description = description
	return c
}

// GetDescription returns the descripiton for the async task
func (c *Config) GetDescription() string {
	return c.description
}

// WithSchedule sets the (gae cron formatted) schedule for the async task
func (c *Config) WithSchedule(schedule string) *Config {
	c.schedule = schedule
	return c
}

// GetSchedule returns the (gae cron formatted) schedule string for the async task
func (c *Config) GetSchedule() string {
	return c.schedule
}

// GetStatus returns a *TaskStatus for the given taskID
func (c *Config) GetStatus(ctx context.Context, taskID string) *TaskStatus {
	t := DefaultAsyncTaskKind.MustGet(ctx, taskID)
	if t == nil {
		return nil
	}
	return t.GetStatus()
}

// GetRecentTasks returns a list of recent tasks ordered by StartAt
func (c *Config) GetRecentTasks(ctx context.Context, n int) []*TaskStatus {
	const defaultNum = 5
	const maxNum = 20
	if n <= 0 || n > maxNum {
		n = defaultNum
	}
	tasks := NewAsyncTaskQuery().Eq("ConfigKey", lazy.New(c.key)).Desc("StartAt").Limit(lazy.New(n)).MustGetAllValues(ctx)
	list := make([]*TaskStatus, len(tasks))
	for i, t := range tasks {
		list[i] = t.GetStatus()
	}
	return list
}

// errors
var (
	ErrAlreadyExists     = fmt.Errorf("already exists")
	ErrNoTaskInstance    = fmt.Errorf("no task instance")
	ErrAlreadyProcessing = fmt.Errorf("already processing")
	ErrPushTaskFailed    = fmt.Errorf("pushtask failed")
)

const exceutionTimeLimit = 10 * time.Minute
const executionTimeWarningThreshold = exceutionTimeLimit * 90 / 100 // 90% of execution limit consumed

// Process run a logic and update AsyncTask record.
func (c *Config) Process(ctx context.Context, taskID string, instancePath string, params url.Values) (*Progress, error) {
	// TODO: There would be orphan tasks that need to be cleaned up since
	// a task progress is tracked on datastore, and we use MustPut() to update the record without any retries.
	t := DefaultAsyncTaskKind.MustGet(ctx, taskID)
	if t == nil {
		return nil, ErrNoTaskInstance
	}
	ctx, logger := xlog.WithContextAndKey(ctx, t.GetLogPrefix(), LoggerKey)
	if t.Status != StatusReady && t.Status != StatusRunning {
		// return 200 for GAE not to retry the task.
		return nil, ErrAlreadyProcessing
	}
	if t.Status == StatusReady {
		logger.Infof("Starting")
		t.StartAt = xtime.Now()
		t.Status = StatusRunning
		DefaultAsyncTaskKind.MustPut(ctx, t)
	}

	var err error
	var progress *Progress

	timeTaken := xtime.Benchmark(func() {
		defer func() {
			var ok bool
			if x := recover(); x != nil {
				err, ok = x.(error)
				if !ok {
					err = fmt.Errorf("%v", x)
				}
				logger.Fatalf("Rescue from panic: %v", err)
			}
		}()
		progress, err = c.logic.Run(ctx, keyvalue.NewQueryProxy(params), t)
	})
	logger.Infof("The logic takes %s", timeTaken)
	if timeTaken > executionTimeWarningThreshold {
		logger.Warnf("The logic takes too long time with 10 minutes limitation, please check and reduce the time in a single execution.")
	}

	if progress != nil {
		t.Progress = append(t.Progress, *progress)
		DefaultAsyncTaskKind.MustPut(ctx, t)
		if err = c.pushTask(ctx, instancePath, progress.Next); err == nil {
			return progress, nil
		}
		// if PushTask fails, AsyncTask is marked as 'failure'
		logger.Infof("PushTask fails due to %v, stopping the task", err)
	}

	// finished the task
	t.FinishAt = xtime.Now()
	t.Progress = nil
	if err == nil {
		t.Status = StatusSuccess
	} else {
		t.Status = StatusFailure
		t.Error = err.Error()
	}
	if err != nil {
		logger.Errorf("Finished with %s in %s: %v", t.Status, t.FinishAt.Sub(t.StartAt), err)
	} else {
		logger.Infof("Finished with %s in %s.", t.Status, t.FinishAt.Sub(t.StartAt))
	}
	DefaultAsyncTaskKind.MustPut(ctx, t)
	return nil, nil
}

// Prepare prepares a AsyncTask record and returns the associated task id.
func (c *Config) Prepare(ctx context.Context, taskID string, instancePath string, params url.Values) (*TaskStatus, error) {
	if c.GetStatus(ctx, taskID) != nil {
		return nil, ErrAlreadyExists
	}
	t := &AsyncTask{}
	t.ID = taskID
	t.ConfigKey = c.key
	if params != nil {
		t.Params = params.Encode()
	}
	t.Status = StatusReady
	t.TaskStore = nil
	ctx, logger := xlog.WithContextAndKey(ctx, t.GetLogPrefix(), LoggerKey)
	DefaultAsyncTaskKind.MustPut(ctx, t)
	if err := c.pushTask(ctx, instancePath, params); err != nil {
		logger.Infof("PushTask fails due to %v, clean up...", err)
		DefaultAsyncTaskKind.MustDelete(ctx, taskID)
		return nil, ErrPushTaskFailed
	}
	logger.Infof("Prepared")
	return t.GetStatus(), nil
}

func (c *Config) pushTask(ctx context.Context, instancePath string, params url.Values) error {
	if params != nil {
		instancePath = fmt.Sprintf("%s?%s", instancePath, params.Encode())
	}
	return c.queue.PushTask(ctx, instancePath, nil)
}
