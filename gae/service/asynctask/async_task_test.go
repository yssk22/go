package asynctask

import (
	"os"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/gae/service"
	"github.com/speedland/go/web"
	"github.com/speedland/go/x/xlog"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestNew(t *testing.T) {
	a := assert.New(t)
	s := service.New("foo")
	cfg := New(s, "/async/task/")
	a.EqStr("foo-async-task", cfg.Queue.Name)
}

func TestLogic(t *testing.T) {
	xlog.SetKeyFilter(LoggerKey, xlog.LevelDebug)
	a := assert.New(t)
	s := service.New("foo")
	cfg := New(s, "/async/task/")
	var processed = false
	var queryValue = ""
	cfg.Implement(Func(func(req *web.Request, task *AsyncTask) (*Progress, error) {
		processed = true
		queryValue = req.Query.GetStringOr("q", "")
		return nil, nil
	}))

	runner := NewTestRunner(t, s)
	task := runner.Run(gaetest.NewContext(), "/foo/async/task/", nil, cfg.Queue.Name)
	a.NotNil(task)
	a.EqInt(int(StatusSuccess), int(task.Status))
}
