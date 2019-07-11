package asynctask

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/yssk22/go/gae/taskqueue"
	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/uuid"
	"github.com/yssk22/go/x/xlog"
	"github.com/yssk22/go/x/xruntime"
	"github.com/yssk22/go/x/xtime"
)

var defaultLogic = LogicFunc(func(ctx context.Context, params *keyvalue.GetProxy, t *AsyncTask) (*Progress, error) {
	_, logger := xlog.WithContextAndKey(ctx, t.GetLogPrefix(), LoggerKey)
	logger.Warnf("no task logic is implemented")
	return nil, nil
})

// Config is an configuration object to define AsyncTask endpoints
type Config struct {
	path        string // task base path.
	queue       *taskqueue.PushQueue
	logic       Logic
	schedule    string
	description string
	timeout     time.Duration
}

const defaultTimeout = 10 * time.Minute

// NewConfig creats a new task configuration
func NewConfig(path string, options ...Option) *Config {
	if !strings.HasSuffix(path, "/") {
		panic(fmt.Errorf("AsyncTask path must ends with '/' (got %q)", path))
	}
	caller := xruntime.CaptureStack(3)[2]
	desc := fmt.Sprintf("defined in %s:%d", caller.FullFilePath, caller.LineNumber)
	c := &Config{
		path:        path,
		logic:       defaultLogic,
		description: desc,
		schedule:    "",
		queue:       taskqueue.DefaultPushQueue,
		timeout:     defaultTimeout,
	}
	for _, opts := range options {
		c, _ = opts(c)
	}
	return c
}

// Option is a func struct to configure
type Option func(c *Config) (*Config, error)

// Queue configures the queue which the task use
func Queue(q *taskqueue.PushQueue) Option {
	return func(c *Config) (*Config, error) {
		c.queue = q
		return c, nil
	}
}

// Description sets the description
func Description(desc string) Option {
	return func(c *Config) (*Config, error) {
		c.description = desc
		return c, nil
	}
}

// Implement sets the logic implementation
func Implement(l Logic) Option {
	return func(c *Config) (*Config, error) {
		c.logic = l
		return c, nil
	}
}

// Func sets the logic implementation with a func. This is a shorthand for asynctask.Implement(asyntask.LogicFunc(f))
func Func(f func(context.Context, *keyvalue.GetProxy, *AsyncTask) (*Progress, error)) Option {
	return func(c *Config) (*Config, error) {
		c.logic = LogicFunc(f)
		return c, nil
	}
}

// Schedule sets the task schedule
func Schedule(sched string) Option {
	return func(c *Config) (*Config, error) {
		c.schedule = sched
		return c, nil
	}
}

// Timeout sets the task timeout
func Timeout(timeout time.Duration) Option {
	return func(c *Config) (*Config, error) {
		c.timeout = timeout
		return c, nil
	}
}

// GetDescription returns the descripiton for the async task
func (c *Config) GetDescription() string {
	return c.description
}

// GetSchedule returns the (gae cron formatted) schedule string for the async task
func (c *Config) GetSchedule() string {
	return c.schedule
}

// GetStatus returns a *TaskStatus for the given taskID
func (c *Config) GetStatus(ctx context.Context, taskID string) *TaskStatus {
	_, t := NewAsyncTaskKind().MustGet(ctx, taskID)
	if t == nil {
		return nil
	}
	if t.ConfigKey != c.path {
		return nil
	}
	return t.GetStatus(c.timeout)
}

// GetRecentTasks returns a list of recent tasks ordered by StartAt
func (c *Config) GetRecentTasks(ctx context.Context, n int) []*TaskStatus {
	const defaultNum = 5
	const maxNum = 20
	if n <= 0 || n > maxNum {
		n = defaultNum
	}
	_, tasks := NewAsyncTaskQuery().EqConfigKey(c.path).DescStartAt().Limit(n).MustGetAll(ctx)
	list := make([]*TaskStatus, len(tasks))
	for i, t := range tasks {
		list[i] = t.GetStatus(c.timeout)
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
	kind := NewAsyncTaskKind()
	_, t := kind.MustGet(ctx, taskID)
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
		kind.MustPut(ctx, t)
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
		kind.MustPut(ctx, t)
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
	kind.MustPut(ctx, t)
	return nil, nil
}

// Prepare prepares a AsyncTask record and returns the associated task id.
func (c *Config) Prepare(ctx context.Context, params url.Values) (*TaskStatus, error) {
	taskID := uuid.New().String()
	instancePath := fmt.Sprintf("%s%s.json", c.path, taskID)
	kind := NewAsyncTaskKind()
	t := &AsyncTask{}
	t.ID = taskID
	t.ConfigKey = c.path
	if params != nil {
		t.Params = params.Encode()
	}
	t.Status = StatusReady
	t.TaskStore = nil
	ctx, logger := xlog.WithContextAndKey(ctx, t.GetLogPrefix(), LoggerKey)
	kind.MustPut(ctx, t)
	if err := c.pushTask(ctx, instancePath, params); err != nil {
		logger.Infof("PushTask fails due to %v, clean up...", err)
		kind.MustDelete(ctx, taskID)
		return nil, ErrPushTaskFailed
	}
	logger.Infof("Prepared")
	return t.GetStatus(c.timeout), nil
}

func (c *Config) pushTask(ctx context.Context, instancePath string, params url.Values) error {
	if params != nil {
		instancePath = fmt.Sprintf("%s?%s", instancePath, params.Encode())
	}
	return c.queue.PushTask(ctx, instancePath, nil)
}
