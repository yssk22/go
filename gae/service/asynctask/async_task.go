// Package asynctask provides async task execution support on GAE apps
package asynctask

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/speedland/go/gae/service"
	"github.com/speedland/go/gae/taskqueue"
	"github.com/speedland/go/uuid"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xcontext"
	"github.com/speedland/go/x/xlog"
	"github.com/speedland/go/x/xtime"
)

// LoggerKey is a xlog key for this package
const LoggerKey = "gae.service.asynctask"

// TaskIDContextKey is a context key for AsyncTaskID
var TaskIDContextKey = xcontext.NewKey("taskid")

// Status is a value to represent the task status
//go:generate enum -type=Status
type Status int

// Available values of Status
const (
	StatusUnknown Status = iota
	StatusReady
	StatusRunning
	StatusSuccess
	StatusFailure
)

// AsyncTask is a record to track a task progress
//go:generate ent -type=AsyncTask
type AsyncTask struct {
	ID        string     `json:"id" ent:"id"`
	Path      string     `json:"path"`
	Query     string     `json:"query"  datastore:",noindex"`
	Status    Status     `json:"status"`
	Error     string     `json:"error" datastore:",noindex"`
	Progress  []Progress `json:"progress" datastore:",noindex"`
	StartAt   time.Time  `json:"start_at"`
	FinishAt  time.Time  `json:"finish_at"`
	UpdatedAt time.Time  `json:"updated_at" ent:"timestamp"`
}

// LastProgress returns the last progress of the task
func (t *AsyncTask) LastProgress() *Progress {
	l := len(t.Progress)
	if l == 0 {
		return nil
	}
	return &t.Progress[l-1]
}

// NewMonitorResponse returns a new *MonitorResponse exposed externally
func (t *AsyncTask) NewMonitorResponse() *MonitorResponse {
	m := &MonitorResponse{
		ID:     t.ID,
		Status: t.Status,
	}
	if !t.StartAt.IsZero() {
		m.StartAt = &(t.StartAt)
	}
	if !t.FinishAt.IsZero() {
		m.FinishAt = &(t.FinishAt)
	}
	if t.Error != "" {
		m.Error = &(t.Error)
	}
	m.Progress = t.LastProgress()
	return m
}

// Progress is a struct that represents the task progress
type Progress struct {
	Total   int        `json:"total,omitempty"`
	Current int        `json:"current,omitempty"`
	Message string     `json:"message,omitempty"`
	Next    url.Values `json:"-" datastore:"-"`
}

// Logic is an interface to execute a task
type Logic interface {
	Run(*web.Request, *AsyncTask) (*Progress, error)
}

// Func is an function to implement Logic
type Func func(*web.Request, *AsyncTask) (*Progress, error)

// Run implements Logic#Run
func (f Func) Run(req *web.Request, t *AsyncTask) (*Progress, error) {
	return f(req, t)
}

// Config is an configuration object to configure endpoints on the task.
type Config struct {
	Queue     *taskqueue.PushQueue
	service   *service.Service
	validator web.Handler
	path      string
	logic     Logic
}

// Implement defines the task logic
func (c *Config) Implement(t Logic) {
	c.logic = t
}

// Schedule adds the cron endpoint for the async task
func (c *Config) Schedule(sched string, description string) {
	p := path.Join(c.path, "cron/")
	c.service.AddCron(p, sched, description,
		web.HandlerFunc(func(req *web.Request, _ web.NextHandler) *response.Response {
			err := c.Queue.PushTask(req.Context(), c.service.Path(p), url.Values{})
			if err != nil {
				return response.NewError(err)
			}
			return response.NewText("OK")
		}),
	)
}

type TriggerResponse struct {
	ID string `json:"id"`
}

type MonitorResponse struct {
	ID       string     `json:"id"`
	Status   Status     `json:"status"`
	StartAt  *time.Time `json:"start_at,omitempty"`
	FinishAt *time.Time `json:"finish_at,omitempty"`
	Error    *string    `json:"error,omitempty"`
	Progress *Progress  `json:"progress,omitempty"`
}

// New adds a new push queue for the asynchronous task execution and returns a *Config value
// to implement the business logic
func New(s *service.Service, path string) *Config {
	if !strings.HasSuffix(path, "/") {
		panic(fmt.Errorf("asynctask path must ends with '/' (got %q)", path))
	}

	name := path[1 : len(path)-1]              // remove prefix/suffix slash.
	name = strings.Replace(name, ":", "", -1)  // remove path parameters (/:param/ => /param/)
	name = strings.Replace(name, "/", "-", -1) // replace '/' with '-' (/path/to/queue/ => path-to-queue)
	queue := s.AddPushQueue(name)

	config := &Config{
		Queue:   queue,
		service: s,
		path:    path,
	}

	// GET /path/:taskid.json
	// endpoint to get the latest status of the given taskid.
	s.Get(fmt.Sprintf("%s:taskid.json", path),
		web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
			t := DefaultAsyncTaskKind.MustGet(req.Context(), req.Params.GetStringOr("taskid", ""))
			if t == nil {
				return nil
			}
			return response.NewJSON(t.NewMonitorResponse())
		}))

	// POST /path/:taskid.json
	// endpoint to execute a task logic, only called via pushtask
	s.Post(fmt.Sprintf("%s:taskid.json", path),
		queue.RequestValidator(),
		web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
			// TODO: There would be orphan tasks that need to be cleaned up since
			// a task progress is tracked on datastore, and we use MustPut() to update the record without any retries.
			t := DefaultAsyncTaskKind.MustGet(req.Context(), req.Params.GetStringOr("taskid", ""))
			if t == nil {
				return nil
			}
			if t.Status != StatusReady && t.Status != StatusRunning {
				return response.NewErrorWithStatus(
					fmt.Errorf("task %q is already in %s", t.ID, t.Status),
					response.HTTPStatusPreconditionFailed,
				)
			}
			logger := xlog.WithContext(context.WithValue(req.Context(), TaskIDContextKey, t.ID)).WithKey(LoggerKey)
			if t.Status == StatusReady {
				logger.Infof("Start a task")
				t.StartAt = xtime.Now()
				t.Status = StatusRunning
				DefaultAsyncTaskKind.MustPut(req.Context(), t)
			}

			var err error
			var progress *Progress
			var resp *response.Response

			func() {
				defer func() {
					var ok bool
					if x := recover(); x != nil {
						err, ok = x.(error)
						if !ok {
							err = fmt.Errorf("%v", x)
						}
					}
				}()
				progress, err = config.logic.Run(req, t)
			}()

			if progress != nil {
				t.Progress = append(t.Progress, *progress)
				if progress.Next != nil {
					DefaultAsyncTaskKind.MustPut(req.Context(), t)
					logger.Infof("The task logic returns a progress with next params, calling a task recursively....")
					err = queue.PushTask(req.Context(), fmt.Sprintf("%s?%s", req.URL.EscapedPath(), progress.Next.Encode()), nil)
					if err == nil {
						// response next parameters to client.
						// this is used only for TestRunner to call the next
						return response.NewJSON(progress.Next)
					}
				} else {
					logger.Infof("The task logic returns a progress without next parameters.")
				}
			} else {
				logger.Infof("The task logic doesn't return a progress.")
			}

			// finished the task
			t.FinishAt = xtime.Now()
			t.Progress = nil
			if err == nil {
				t.Status = StatusSuccess
				resp = response.NewJSON(true)
			} else {
				t.Error = err.Error()
				t.Status = StatusFailure
				resp = response.NewError(err)
			}
			DefaultAsyncTaskKind.MustPut(req.Context(), t)
			logger.Infof("The task finished with %s(tt=%s).", t.Status, t.FinishAt.Sub(t.StartAt))
			return resp
		}))

	// POST /path/
	// endpoint to create a new task record and call /path/:taskid.json via pushtask.
	s.Post(path,
		web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
			t := &AsyncTask{}
			t.ID = uuid.New().String()
			t.Path = req.URL.Path
			t.Query = req.URL.RawQuery
			t.Status = StatusReady
			DefaultAsyncTaskKind.MustPut(req.Context(), t)
			taskPath := fmt.Sprintf("%s%s.json?%s", req.URL.Path, t.ID, req.URL.Query().Encode())
			if err := queue.PushTask(req.Context(), taskPath, nil); err != nil {
				panic(err)
			}
			logger := xlog.WithContext(context.WithValue(req.Context(), TaskIDContextKey, t.ID)).WithKey(LoggerKey)
			logger.Infof("An AsyncTask created: %s (path:%s, queue:%s)", t.ID, taskPath, queue.Name)
			return response.NewJSONWithStatus(
				&TriggerResponse{t.ID},
				response.HTTPStatusCreated,
			)
		}))
	return config
}
