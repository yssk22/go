package asynctaskrunner

import (
	"net/url"
	"os"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"

	"github.com/yssk22/go/gae/gaetest"
	"github.com/yssk22/go/gae/service"
	"github.com/yssk22/go/gae/service/asynctask"
	"github.com/yssk22/go/keyvalue"
	"context"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestLogic(t *testing.T) {
	a := assert.New(t)
	s := service.New("foo")
	var processed = false
	var queryValue = ""
	queue := s.AddPushQueue("myqueue")
	s.AsyncTask("/async/task/", asynctask.NewConfig(
		"async-task",
		s.AddPushQueue("myqueue"),
		asynctask.Func(func(ctx context.Context, params *keyvalue.GetProxy, t *asynctask.AsyncTask) (*asynctask.Progress, error) {
			processed = true
			queryValue = params.GetStringOr("q", "")
			return nil, nil
		}),
	))

	runner := NewAsyncTaskRunner(t, s)
	task := runner.Run(gaetest.NewContext(), "/foo/async/task/", nil, queue.Name)
	a.NotNil(task)
	a.EqInt(int(asynctask.StatusSuccess), int(task.Status))
}

func TestTaskStore(t *testing.T) {
	a := assert.New(t)
	s := service.New("foo")
	queue := s.AddPushQueue("myqueue")
	type V struct {
		Foo string
	}
	var v V
	var queryValue = ""
	s.AsyncTask("/async/task/", asynctask.NewConfig(
		"async-task",
		s.AddPushQueue("myqueue"),
		asynctask.Func(func(ctx context.Context, params *keyvalue.GetProxy, t *asynctask.AsyncTask) (*asynctask.Progress, error) {
			if t.IsStoreEmpty() {
				t.SaveStore(&V{
					Foo: "bar",
				})
				return &asynctask.Progress{
					Total:   2,
					Current: 1,
					Next:    url.Values{"a": []string{"1"}},
				}, nil
			}
			t.LoadStore(&v)
			queryValue = params.GetStringOr("q", "")
			return nil, nil
		}),
	))

	runner := NewAsyncTaskRunner(t, s)
	task := runner.Run(gaetest.NewContext(), "/foo/async/task/", nil, queue.Name)
	a.NotNil(task)
	a.EqStr("bar", v.Foo)
	a.EqInt(int(asynctask.StatusSuccess), int(task.Status))
}
