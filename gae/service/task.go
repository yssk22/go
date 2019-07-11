package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/yssk22/go/gae/service/apierrors"
	"github.com/yssk22/go/gae/service/asynctask"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
	"github.com/yssk22/go/x/xerrors"
)

// Task is a defined task in the service.
type Task struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
	config      *asynctask.Config
}

// AsyncTask defines endpoints for asynctask execution
func (s *Service) AsyncTask(path string, options ...asynctask.Option) {
	fullPath := s.Path(path)
	instancePathTemplate := fmt.Sprintf("%s:taskid.json", path)
	taskConfig := asynctask.NewConfig(fullPath, options...)

	// Create a new task instasnce
	s.Post(path, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		status, err := taskConfig.Prepare(req.Context(), req.Request.URL.Query())
		xerrors.MustNil(err)
		return response.NewJSONWithStatus(status, response.HTTPStatusCreated)
	}))

	// Get the instance status
	s.Get(instancePathTemplate, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		status := taskConfig.GetStatus(req.Context(), req.Params.GetStringOr("taskid", ""))
		if status == nil {
			return nil
		}
		return response.NewJSON(status)
	}))

	const TaskQueueHeader = "X-AppEngine-TaskName"
	// Run a task instance
	s.Post(instancePathTemplate, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		if req.Header.Get(TaskQueueHeader) == "" {
			return apierrors.Forbidden.ToResponse()
		}
		taskID := req.Params.GetStringOr("taskid", "")
		progress, err := taskConfig.Process(req.Context(), taskID, fmt.Sprintf("%s%s.json", fullPath, taskID), req.Request.URL.Query())
		if err != nil {
			if err == asynctask.ErrNoTaskInstance {
				return nil
			}
			return (&apierrors.Error{
				Code:    "invalid_asynctask_call",
				Message: err.Error(),
				Status:  response.HTTPStatusBadRequest,
			}).ToResponse()
		}
		if progress == nil {
			return response.NewJSON(true)
		}
		return response.NewJSON(progress.Next)
	}))

	// Get Recent Tasks
	s.Get(path, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return response.NewJSON(taskConfig.GetRecentTasks(req.Context(), req.Query.GetIntOr("n", 5)))
	}))

	if schedule := taskConfig.GetSchedule(); schedule != "" {
		const CronHeader = "X-AppEngine-Cron"
		s.AddCron(fmt.Sprintf("%s/cron/", path), schedule, taskConfig.GetDescription(), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
			if req.Header.Get(CronHeader) == "" {
				return apierrors.Forbidden.ToResponse()
			}
			params := req.Request.URL.Query()
			params.Set("cron", "true")
			status, err := taskConfig.Prepare(req.Context(), params)
			xerrors.MustNil(err)
			return response.NewJSONWithStatus(status, response.HTTPStatusCreated)
		}))
	}
	s.tasks = append(s.tasks, &Task{
		ID:          s.Path(path),
		Path:        s.Path(path),
		Description: taskConfig.GetDescription(),
		Schedule:    taskConfig.GetSchedule(),
		config:      taskConfig,
	})
}

// GetTasks returns a list of queues defined in the service
func (s *Service) GetTasks() []*Task {
	return s.tasks
}

// RunTask runs the task specified by the path adhocly
func (s *Service) RunTask(ctx context.Context, path string, params url.Values) (*asynctask.TaskStatus, error) {
	fullPath := s.Path(path)
	for _, t := range s.tasks {
		if t.Path == fullPath {
			if params == nil {
				params = url.Values{}
			}
			params.Set("__run_task", "true")
			return t.config.Prepare(ctx, params)
		}
	}
	return nil, fmt.Errorf("no task path found at %s", fullPath)
}
