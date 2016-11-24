package service

import (
	"testing"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/web/httptest"
	"github.com/speedland/go/web/response"
)

func TestService_serveQueue(t *testing.T) {
	a := httptest.NewAssert(t)
	s := New("test")
	recorder := gaetest.NewRecorder(s)
	resp := recorder.TestGet("/test/__/queue.yaml")
	a.Status(response.HTTPStatusOK, resp)
}
