package builtin

import (
	"fmt"
	"path"

	"github.com/yssk22/go/gae/service"
	"github.com/yssk22/go/gae/service/apierrors"
	"github.com/yssk22/go/gae/service/auth"
	"github.com/yssk22/go/services/facebook/messenger"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
	"github.com/yssk22/go/x/xerrors"

	"github.com/yssk22/go/x/xlog"
	"google.golang.org/appengine"
)

func setupConfigAPIs(s *service.Service) {
	const basePath = "/admin/api/configs/"
	c := s.Config
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
	const basePath = "/admin/api/tasks/"
	s.Get(basePath, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return response.NewJSON(s.GetTasks())
	}))
}

func setupAuthAPIs(s *service.Service) {
	const basePath = "/auth/"
	s.Get(path.Join(basePath, "me.json"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx, err := appengine.Namespace(req.Context(), s.APIConfig.AuthNamespace)
		xerrors.MustNil(err)
		a, err := auth.GetCurrent(ctx)
		if err != nil {
			return apierrors.ServerError.ToResponse()
		}
		return response.NewJSON(a)
	}))

	s.Post(path.Join(basePath, "login/facebook/"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		token := req.Form.GetStringOr("access_token", "")
		ctx, err := appengine.Namespace(req.Context(), s.APIConfig.AuthNamespace)
		xerrors.MustNil(err)
		ctx, logger := xlog.WithContext(ctx, "")
		a, err := auth.Facebook(ctx, s.Config.NewHTTPClient(ctx), token)
		if err != nil {
			logger.Infof("failed to authenticate facebook with %q: %v", token, err)
			return apierrors.BadRequest.ToResponse()
		}
		if err = auth.SetCurrent(ctx, a); err != nil {
			panic(err)
		}
		return response.NewJSON(a)
	}))

	s.Get(path.Join(basePath, "logout/"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		if !appengine.IsDevAppServer() {
			// available only in DevAppServer
			return nil
		}
		auth.DeleteCurrent(req.Context())
		return response.NewJSON("OK")
	}))

	s.Post(path.Join(basePath, "logout/"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		auth.DeleteCurrent(req.Context())
		return response.NewJSON("OK")
	}))
}

func setupWebhooks(s *service.Service) {
	webhook := s.APIConfig.MessengerWebHook
	if webhook == nil {
		return
	}
	s.Get(path.Join("__webhook/messenger/"), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		token := s.Config.GetMessengerVerificationToken(req.Context())
		if token == "" {
			panic(fmt.Errorf("no messenger verification token is configured"))
		}
		return messenger.NewVericationHandler(token).Process(req, next)
	}))
	s.Post(path.Join("__webhook/messenger/"), messenger.NewWebhookHandler(webhook))
}
