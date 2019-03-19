package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/yssk22/go/gae/service/apierrors"
	"github.com/yssk22/go/gae/service/asynctask"
	"github.com/yssk22/go/uuid"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
	"github.com/yssk22/go/x/xerrors"
)

// Task is a defined task in the service.
// @flow
type Task struct {
	Path        string `json:"path"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
	config      *asynctask.Config
}

// AsyncTask defines endpoints for asynctask execution
func (s *Service) AsyncTask(path string, taskConfig *asynctask.Config) {
	if !strings.HasSuffix(path, "/") {
		panic(fmt.Errorf("AsyncTask path must ends with '/' (got %q)", path))
	}
	fullPath := s.Path(path)
	s.Get(fmt.Sprintf("%s:taskid.json", path), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		status := taskConfig.GetStatus(req.Context(), req.Params.GetStringOr("taskid", ""))
		if status == nil {
			return nil
		}
		return response.NewJSON(status)
	}))

	const TaskQueueHeader = "X-AppEngine-TaskName"
	s.Post(fmt.Sprintf("%s:taskid.json", path), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
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

	s.Post(path, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		taskID := uuid.New().String()
		status, err := taskConfig.Prepare(req.Context(), taskID, fmt.Sprintf("%s%s.json", fullPath, taskID), req.Request.URL.Query())
		xerrors.MustNil(err)
		return response.NewJSONWithStatus(status, response.HTTPStatusCreated)
	}))

	s.Get(path, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return response.NewJSON(taskConfig.GetRecentTasks(req.Context(), req.Query.GetIntOr("n", 5)))
	}))

	if schedule := taskConfig.GetSchedule(); schedule != "" {
		const CronHeader = "X-AppEngine-Cron"
		s.AddCron(fmt.Sprintf("%s/cron/", path), schedule, taskConfig.GetDescription(), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
			if req.Header.Get(CronHeader) == "" {
				return apierrors.Forbidden.ToResponse()
			}
			taskID := uuid.New().String()
			params := req.Request.URL.Query()
			params.Set("cron", "true")
			status, err := taskConfig.Prepare(req.Context(), taskID, fmt.Sprintf("%s%s.json", fullPath, taskID), params)
			xerrors.MustNil(err)
			return response.NewJSONWithStatus(status, response.HTTPStatusCreated)
		}))
	}
	s.tasks = append(s.tasks, &Task{
		Path:        s.Path(path),
		Key:         taskConfig.GetKey(),
		Description: taskConfig.GetDescription(),
		Schedule:    taskConfig.GetSchedule(),
		config:      taskConfig,
	})
}

// GetTasks returns a list of queues defined in the service
func (s *Service) GetTasks() []*Task {
	return s.tasks
}

func (s *Service) RunTask(ctx context.Context, path string, params url.Values) (*asynctask.TaskStatus, error) {
	fullPath := s.Path(path)
	for _, t := range s.tasks {
		if t.Path == fullPath {
			taskID := uuid.New().String()
			if params == nil {
				params = url.Values{}
			}
			params.Set("__run_task", "true")
			return t.config.Prepare(ctx, taskID, fmt.Sprintf("%s%s.json", fullPath, taskID), params)
		}
	}
	return nil, fmt.Errorf("no task path found at %s", fullPath)
}
