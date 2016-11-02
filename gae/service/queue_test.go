package service

import (
	"testing"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/httptest"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xlog"
)

func TestService_serveQueue(t *testing.T) {
	xlog.SetKeyFilter(web.LoggerKeyRouter, xlog.LevelDebug)
	a := httptest.NewAssert(t)
	s := New("test")
	recorder := gaetest.NewRecorder(s)
	resp := recorder.TestGet("/test/__/queue.yaml")
	a.Status(response.HTTPStatusOK, resp)
}
