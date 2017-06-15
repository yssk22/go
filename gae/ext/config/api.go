package config

import (
	"path"

	"github.com/speedland/go/gae/service"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
)

// SetupAPI sets up API endpoints for getting/updating configs by
//
//    - GET /{basePath}/
//    - GET /{basePath}/:key.json
//    - PUT /{basePath}/:key.json
//
func SetupAPI(s *service.Service, c *Config, basePath string) {
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
