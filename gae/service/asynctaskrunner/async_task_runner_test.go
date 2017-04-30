package asynctaskrunner

import (
	"net/url"
	"os"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/gae/service"
	"github.com/speedland/go/gae/service/asynctask"
	"github.com/speedland/go/web"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestNew(t *testing.T) {
	a := assert.New(t)
	s := service.New("foo")
	cfg := asynctask.New(s, "/async/task/")
	a.EqStr("foo-async-task", cfg.Queue.Name)
}

func TestLogic(t *testing.T) {
	a := assert.New(t)
	s := service.New("foo")
	cfg := asynctask.New(s, "/async/task/")
	var processed = false
	var queryValue = ""
	cfg.Implement(asynctask.Func(func(req *web.Request, task *asynctask.AsyncTask) (*asynctask.Progress, error) {
		processed = true
		queryValue = req.Query.GetStringOr("q", "")
		return nil, nil
	}))

	runner := NewAsyncTaskRunner(t, s)
	task := runner.Run(gaetest.NewContext(), "/foo/async/task/", nil, cfg.Queue.Name)
	a.NotNil(task)
	a.EqInt(int(asynctask.StatusSuccess), int(task.Status))
}

func TestTaskStore(t *testing.T) {
	a := assert.New(t)
	s := service.New("foo")
	cfg := asynctask.New(s, "/async/task/")
	var queryValue = ""
	type V struct {
		Foo string
	}
	var v V
	cfg.Implement(asynctask.Func(func(req *web.Request, task *asynctask.AsyncTask) (*asynctask.Progress, error) {
		if task.IsStoreEmpty() {
			task.SaveStore(&V{
				Foo: "bar",
			})
			return &asynctask.Progress{
				Total:   2,
				Current: 1,
				Next:    url.Values{"a": []string{"1"}},
			}, nil
		}
		task.LoadStore(&v)
		queryValue = req.Query.GetStringOr("q", "")
		return nil, nil
	}))

	runner := NewAsyncTaskRunner(t, s)
	task := runner.Run(gaetest.NewContext(), "/foo/async/task/", nil, cfg.Queue.Name)
	a.NotNil(task)
	a.EqStr("bar", v.Foo)
	a.EqInt(int(asynctask.StatusSuccess), int(task.Status))
}
