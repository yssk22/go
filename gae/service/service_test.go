package service

import (
	"os"
	"testing"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/httptest"
	"github.com/speedland/go/web/response"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestService(t *testing.T) {
	a := httptest.NewAssert(t)
	s := New("test")
	s.Get("/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		svc, ok := req.Context().Value(ContextKey).(*Service)
		if ok {
			return response.NewText(svc.Key())
		}
		return next(req)
	}))
	recorder := gaetest.NewRecorder(s)
	resp := recorder.TestGet("/test/")
	a.Status(response.HTTPStatusOK, resp)
	a.Body("test", resp)
}
