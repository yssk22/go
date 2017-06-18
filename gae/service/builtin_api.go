package service

import (
	"path"

	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
)

// BuiltInAPIConfig is a configuration object for ActivateBuiltInAPIs, which actiate the following endpoints
//
// [config]
//
//    - GET /{ConfigAPIBasePath}/
//    - GET /{ConfigAPIBasePath}/:key.json
//    - PUT /{ConfigAPIBasePath}/:key.json
//
// [asynctask]
//    - GET /{AsyncTaskListAPIPath}/
//
type BuiltInAPIConfig struct {
	ConfigAPIBasePath    string
	AsyncTaskListAPIPath string
}

// ActivateEndpoints sets up builtin API endpoints on the *Service
func (bc *BuiltInAPIConfig) ActivateEndpoints(s *Service) *Service {
	bc.activateConfigAPI(s)
	bc.activateAsyncTaskListAPI(s)
	return s
}

// DefaultBuiltinAPIConfig is a default object of BuiltInAPIConfig
var DefaultBuiltinAPIConfig = &BuiltInAPIConfig{
	ConfigAPIBasePath:    "/admin/api/configs/",
	AsyncTaskListAPIPath: "/admin/api/asynctasks/",
}

func (bc *BuiltInAPIConfig) activateConfigAPI(s *Service) {
	if bc.ConfigAPIBasePath == "" {
		return
	}
	c := s.Config
	basePath := bc.ConfigAPIBasePath
	s.Get(basePath, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return response.NewJSON(c.All(req.Context()))
	}))

	s.Get(path.Join(basePath, ":key.json"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		cfg := c.Get(req.Context(), req.Params.GetStringOr("key", ""))
		if cfg == nil {
			return nil
		}
		return response.NewJSON(cfg)
	}))

	s.Put(path.Join(basePath, ":key.json"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		cfg := c.Get(req.Context(), req.Params.GetStringOr("key", ""))
		if cfg == nil {
			return nil
		}
		cfg.UpdateByForm(req.Form)
		c.Set(req.Context(), cfg)
		return response.NewJSON(cfg)
	}))
}

// AsyncTaskListItem is an list item of AsyncTaskList API reponse.
type AsyncTaskListItem struct {
	Path        string `json:"path"`
	Key         string `json:"key"`
	Description string `json:"description"`
}

func (bc *BuiltInAPIConfig) activateAsyncTaskListAPI(s *Service) {
	if bc.AsyncTaskListAPIPath == "" {
		return
	}
	s.Get(bc.AsyncTaskListAPIPath, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		var list []*AsyncTaskListItem
		for _, def := range s.tasks {
			list = append(list, &AsyncTaskListItem{
				Path:        def.path,
				Key:         def.config.GetKey(),
				Description: def.config.GetDescription(),
			})
		}
		return response.NewJSON(list)
	}))
}
