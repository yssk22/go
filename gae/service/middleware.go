package service

import (
	"bytes"
	"fmt"

	"github.com/speedland/go/gae/service/apierrors"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xlog"
	"context"
	"google.golang.org/appengine"
)

const loggerKeyErrorResponse = "jp.poiku.error"

type httpError interface {
	Status() response.HTTPStatus
}

var namespaceMiddleware = func(s *Service) web.HandlerFunc {
	return web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx, err := appengine.Namespace(req.Context(), s.namespace)
		if err != nil {
			panic(err)
		}
		ctx = context.WithValue(ctx, ContextKey, s)
		return next(req.WithContext(ctx))
	})
}

var errInternalServerError = &apierrors.Error{
	Code:    "internal_server_error",
	Message: "An internal server error occrurred. Please try again later",
	Status:  response.HTTPStatusInternalServerError,
}

var errorMiddleware = web.HandlerFunc(
	func(req *web.Request, next web.NextHandler) *response.Response {
		s := FromContext(req.Context())
		var resp *response.Response
		func() {
			defer func() {
				if x := recover(); x != nil {
					var status = response.HTTPStatusInternalServerError
					var err error
					var ok bool
					err, ok = x.(error)
					if !ok {
						err = fmt.Errorf("%v", err)
					}
					httpe, ok := err.(httpError)
					if ok {
						status = httpe.Status()
					}
					xlog.WithContext(req.Context()).WithKey(loggerKeyErrorResponse).Fatalf(
						"Recovering error from panic: %v", err,
					)
					if s.OnError != nil {
						resp = s.OnError(req, err)
					} else {
						if !appengine.IsDevAppServer() {
							resp = errInternalServerError.ToResponse()
						} else {
							resp = (&apierrors.Error{
								Code:    errInternalServerError.Code,
								Message: err.Error(),
								Status:  status,
							}).ToResponse()
						}
					}
				}
			}()
			resp = next(req)
		}()
		if resp != nil {
			code := int(resp.Status)
			if code >= 500 {
				var buff bytes.Buffer
				resp.Body.Render(req.Context(), &buff)
				xlog.WithContext(req.Context()).WithKey(loggerKeyErrorResponse).Errorf(
					"Internal Server Error: %v", buff.String(),
				)
			} else if code > 404 {
				var buff bytes.Buffer
				resp.Body.Render(req.Context(), &buff)
				xlog.WithContext(req.Context()).WithKey(loggerKeyErrorResponse).Warnf(
					"Unusual client error: %v", buff.String(),
				)
			}
		}
		return resp
	},
)

var initMiddleware = web.HandlerFunc(
	func(req *web.Request, next web.NextHandler) *response.Response {
		s := FromContext(req.Context())
		if s.Init != nil {
			s.once.Do(func() {
				s.Init(req)
			})
		}
		return next(req)
	},
)

var everyMiddleware = web.HandlerFunc(
	func(req *web.Request, next web.NextHandler) *response.Response {
		s := FromContext(req.Context())
		if s.Every != nil {
			s.Every(req)
		}
		return next(req)
	},
)
