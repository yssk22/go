// Package asynctask provides async task execution support on GAE apps
package asynctask

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/speedland/go/x/xcontext"
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

// TaskStore is a type alias for []byte
type TaskStore []byte

// MarshalJSON implements json.MarshalJSON()
func (cs TaskStore) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(cs))
}

// UnmarshalJSON implements json.Unmarshaler#UnmarshalJSON([]byte)
func (cs *TaskStore) UnmarshalJSON(b []byte) error {
	var v []byte
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	*cs = TaskStore(v)
	return nil
}

// AsyncTask is a record to track a task progress
//go:generate ent -type=AsyncTask
type AsyncTask struct {
	ID        string     `json:"id" ent:"id"`
	ConfigKey string     `json:"config_key"`
	Params    string     `json:"params" datastore:",noindex"`
	Status    Status     `json:"status"`
	Error     string     `json:"error" datastore:",noindex"`
	Progress  []Progress `json:"progress" datastore:",noindex"`
	TaskStore TaskStore  `json:"taskstore" datastore",noindex"`
	StartAt   time.Time  `json:"start_at"`
	FinishAt  time.Time  `json:"finish_at"`
	UpdatedAt time.Time  `json:"updated_at" ent:"timestamp"`

	// Deprecated
	Path  string `json:"path"  datastore:",noindex"`
	Query string `json:"query"  datastore:",noindex"`
}

// IsStoreEmpty returns whether TaskStore field is empty or not
func (t *AsyncTask) IsStoreEmpty() bool {
	return t.TaskStore == nil || len(t.TaskStore) == 0
}

// SaveStore updates the task store
func (t *AsyncTask) SaveStore(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	t.TaskStore = TaskStore(b)
	return nil
}

// LoadStore updates the task store
func (t *AsyncTask) LoadStore(v interface{}) error {
	return json.Unmarshal(t.TaskStore, v)
}

// LastProgress returns the last progress of the task
func (t *AsyncTask) LastProgress() *Progress {
	l := len(t.Progress)
	if l == 0 {
		return nil
	}
	return &t.Progress[l-1]
}

// GetStatus returns a new *TaskStatus exposed to clients.
func (t *AsyncTask) GetStatus() *TaskStatus {
	st := &TaskStatus{
		ID:     t.ID,
		Status: t.Status,
	}
	if !t.StartAt.IsZero() {
		st.StartAt = &(t.StartAt)
	}
	if !t.FinishAt.IsZero() {
		st.FinishAt = &(t.FinishAt)
	}
	if t.Error != "" {
		st.Error = &(t.Error)
	}
	st.Progress = t.LastProgress()
	return st
}

// Progress is a struct that represents the task progress
type Progress struct {
	Total   int        `json:"total,omitempty"`
	Current int        `json:"current,omitempty"`
	Message string     `json:"message,omitempty"`
	Next    url.Values `json:"-" datastore:"-"`
}

// TaskStatus is a struct that can be used in task manager clients.
type TaskStatus struct {
	ID       string     `json:"id"`
	Status   Status     `json:"status"`
	StartAt  *time.Time `json:"start_at,omitempty"`
	FinishAt *time.Time `json:"finish_at,omitempty"`
	Error    *string    `json:"error,omitempty"`
	Progress *Progress  `json:"progress,omitempty"`
}

// // Config is an configuration object to configure endpoints on the task.
// type Config struct {
// 	Queue     *taskqueue.PushQueue
// 	service   *service.Service
// 	validator web.Handler
// 	path      string
// 	logic     Logic
// }

// // Implement defines the task logic
// func (c *Config) Implement(t Logic) {
// 	c.logic = t
// }

// // Schedule adds the cron endpoint for the async task
// func (c *Config) Schedule(sched string, description string) {
// 	p := path.Join(c.path, "cron/")
// 	c.service.AddCron(p, sched, description,
// 		web.HandlerFunc(func(req *web.Request, _ web.NextHandler) *response.Response {
// 			path := fmt.Sprintf("%s?cron=true", c.service.Path(c.path))
// 			err := c.Queue.PushTask(req.Context(), path, url.Values{})
// 			if err != nil {
// 				return response.NewError(err)
// 			}
// 			return response.NewText("OK")
// 		}),
// 	)
// }

// type TriggerResponse struct {
// 	ID string `json:"id"`
// }

// type MonitorResponse struct {
// 	ID       string     `json:"id"`
// 	Status   Status     `json:"status"`
// 	StartAt  *time.Time `json:"start_at,omitempty"`
// 	FinishAt *time.Time `json:"finish_at,omitempty"`
// 	Error    *string    `json:"error,omitempty"`
// 	Progress *Progress  `json:"progress,omitempty"`
// }

// // New adds a new push queue for the asynchronous task execution and returns a *Config value
// // to implement the business logic
// func New(s *service.Service, path string) *Config {
// 	if !strings.HasSuffix(path, "/") {
// 		panic(fmt.Errorf("asynctask path must ends with '/' (got %q)", path))
// 	}

// 	name := path[1 : len(path)-1]              // remove prefix/suffix slash.
// 	name = strings.Replace(name, ":", "", -1)  // remove path parameters (/:param/ => /param/)
// 	name = strings.Replace(name, "/", "-", -1) // replace '/' with '-' (/path/to/queue/ => path-to-queue)
// 	queue := s.AddPushQueue(name)

// 	config := &Config{
// 		Queue:   queue,
// 		service: s,
// 		path:    path,
// 	}

// 	// POST /path/:taskid.json
// 	// endpoint to execute a task logic, only called via pushtask
// 	s.Post(fmt.Sprintf("%s:taskid.json", path),
// 		queue.RequestValidator(),
// 		web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
// 			// TODO: There would be orphan tasks that need to be cleaned up since
// 			// a task progress is tracked on datastore, and we use MustPut() to update the record without any retries.
// 			t := DefaultAsyncTaskKind.MustGet(req.Context(), req.Params.GetStringOr("taskid", ""))
// 			if t == nil {
// 				return nil
// 			}
// 			logger := xlog.WithContext(context.WithValue(req.Context(), TaskIDContextKey, t.ID)).WithKey(LoggerKey)
// 			if t.Status != StatusReady && t.Status != StatusRunning {
// 				// return 200 for GAE not to retry the task.
// 				logger.Warnf("task %q is already in %s", t.ID, t.Status)
// 				return response.NewText("OK")
// 			}
// 			if t.Status == StatusReady {
// 				logger.Infof("Start a task")
// 				t.StartAt = xtime.Now()
// 				t.Status = StatusRunning
// 				DefaultAsyncTaskKind.MustPut(req.Context(), t)
// 			}

// 			var err error
// 			var progress *Progress
// 			var resp *response.Response

// 			func() {
// 				defer func() {
// 					var ok bool
// 					if x := recover(); x != nil {
// 						err, ok = x.(error)
// 						if !ok {
// 							err = fmt.Errorf("%v", x)
// 						}
// 						logger.Fatalf("rescue from panic: %v", err)
// 					}
// 				}()
// 				progress, err = config.logic.Run(req, t)
// 			}()

// 			if progress != nil {
// 				t.Progress = append(t.Progress, *progress)
// 				if progress.Next != nil {
// 					DefaultAsyncTaskKind.MustPut(req.Context(), t)
// 					logger.Infof("The task logic returns a progress with next params, calling a task recursively....")
// 					err = queue.PushTask(req.Context(), fmt.Sprintf("%s?%s", req.URL.EscapedPath(), progress.Next.Encode()), nil)
// 					if err == nil {
// 						// response next parameters to client.
// 						// this is used only for TestRunner to call the next
// 						return response.NewJSON(progress.Next)
// 					}
// 				} else {
// 					logger.Infof("The task logic returns a progress without next parameters.")
// 				}
// 			} else {
// 				logger.Infof("The task logic doesn't return a progress.")
// 			}

// 			// finished the task
// 			t.FinishAt = xtime.Now()
// 			t.Progress = nil
// 			if err == nil {
// 				t.Status = StatusSuccess
// 				resp = response.NewJSON(true)
// 			} else {
// 				t.Error = err.Error()
// 				t.Status = StatusFailure
// 				resp = response.NewError(err)
// 			}
// 			DefaultAsyncTaskKind.MustPut(req.Context(), t)
// 			logger.Infof("The task finished with %s(tt=%s).", t.Status, t.FinishAt.Sub(t.StartAt))
// 			return resp
// 		}))

// 	// POST /path/
// 	// endpoint to create a new task record and call /path/:taskid.json via pushtask.
// 	s.Post(path,
// 		web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
// 			t := &AsyncTask{}
// 			t.ID = uuid.New().String()
// 			t.Path = req.URL.Path
// 			t.Query = req.URL.RawQuery
// 			t.Status = StatusReady
// 			t.TaskStore = nil
// 			DefaultAsyncTaskKind.MustPut(req.Context(), t)
// 			taskPath := fmt.Sprintf("%s%s.json?%s", req.URL.Path, t.ID, req.URL.Query().Encode())
// 			if err := queue.PushTask(req.Context(), taskPath, nil); err != nil {
// 				panic(err)
// 			}
// 			logger := xlog.WithContext(context.WithValue(req.Context(), TaskIDContextKey, t.ID)).WithKey(LoggerKey)
// 			logger.Infof("An AsyncTask created: %s (path:%s, queue:%s, cron=%t)", t.ID, taskPath, queue.Name, req.Query.GetStringOr("cron", "") == "true")
// 			return response.NewJSONWithStatus(
// 				&TriggerResponse{t.ID},
// 				response.HTTPStatusCreated,
// 			)
// 		}))
// 	return config
// }
