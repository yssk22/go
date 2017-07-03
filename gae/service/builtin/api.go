package builtin

import (
	"path"

	"github.com/speedland/go/gae/service"
	"github.com/speedland/go/gae/service/apierrors"
	"github.com/speedland/go/gae/service/auth"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xerrors"

	"github.com/speedland/go/x/xlog"
	"google.golang.org/appengine"
)

func setupConfigAPIs(s *service.Service) {
	if s.APIConfig == nil {
		return
	}
	if s.APIConfig.ConfigAPIBasePath == "" {
		return
	}
	c := s.Config
	basePath := s.APIConfig.ConfigAPIBasePath
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

func setupAsyncTaskListAPIs(s *service.Service) {
	if s.APIConfig == nil {
		return
	}
	if s.APIConfig.AsyncTaskListAPIPath == "" {
		return
	}
	s.Get(s.APIConfig.AsyncTaskListAPIPath, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return response.NewJSON(s.GetTasks())
	}))
}

func setupAuthAPIs(s *service.Service) {
	if s.APIConfig == nil {
		return
	}
	if s.APIConfig.AuthAPIBasePath == "" {
		return
	}
	s.Get(path.Join(s.APIConfig.AuthAPIBasePath, "me.json"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx, err := appengine.Namespace(req.Context(), s.APIConfig.AuthNamespace)
		xerrors.MustNil(err)
		a, err := auth.GetCurrent(ctx)
		if err != nil {
			return apierrors.ServerError.ToResponse()
		}
		return response.NewJSON(a)
	}))

	s.Post(path.Join(s.APIConfig.AuthAPIBasePath, "login/facebook/"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		token := req.Form.GetStringOr("access_token", "")
		ctx, err := appengine.Namespace(req.Context(), s.APIConfig.AuthNamespace)
		xerrors.MustNil(err)
		a, err := auth.Facebook(ctx, s.Config.NewHTTPClient(ctx), token)
		if err != nil {
			xlog.WithContext(ctx).Infof("failed to authenticate facebook with %q: %v", token, err)
			return apierrors.BadRequest.ToResponse()
		}
		if err = auth.SetCurrent(ctx, a); err != nil {
			panic(err)
		}
		return response.NewJSON(a)
	}))

	s.Get(path.Join(s.APIConfig.AuthAPIBasePath, "logout/"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		if !appengine.IsDevAppServer() {
			// available only in DevAppServer
			return nil
		}
		auth.DeleteCurrent(req.Context())
		return response.NewJSON("OK")
	}))

	s.Post(path.Join(s.APIConfig.AuthAPIBasePath, "logout/"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		auth.DeleteCurrent(req.Context())
		return response.NewJSON("OK")
	}))
}
