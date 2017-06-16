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
type BuiltInAPIConfig struct {
	ConfigAPIBasePath string
}

// ActivateEndpoints sets up builtin API endpoints on the *Service
func (bc *BuiltInAPIConfig) ActivateEndpoints(s *Service) {
	bc.activateConfigAPI(s)
}

// DefaultBuiltinAPIConfig is a default object of BuiltInAPIConfig
var DefaultBuiltinAPIConfig = &BuiltInAPIConfig{
	ConfigAPIBasePath: "/admin/api/configs/",
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
