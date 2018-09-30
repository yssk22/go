// Package asynctask provides async task execution support on GAE apps
package asynctask

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/yssk22/go/x/xcontext"
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

// GetLogPrefix returns a prefix string for logger
func (t *AsyncTask) GetLogPrefix() string {
	return fmt.Sprintf("[AsyncTask:%s:%s] ", t.ConfigKey, t.ID)
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
