package asynctaskrunner

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"google.golang.org/appengine"

	"golang.org/x/net/context"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/gae/service"
	"github.com/speedland/go/gae/service/asynctask"
	"github.com/speedland/go/web/httptest"
	"github.com/speedland/go/web/response"
)

// AsyncTaskRunner is a runner object to run async tasks in your test.
type AsyncTaskRunner struct {
	t       *testing.T
	service *service.Service
}

// NewAsyncTaskRunner returns a new *AsyncTaskRunner object to run async tasks.
func NewAsyncTaskRunner(t *testing.T, service *service.Service) *AsyncTaskRunner {
	return &AsyncTaskRunner{
		t:       t,
		service: service,
	}
}

// Run executes the task and wait for the completion.
func (runner *AsyncTaskRunner) Run(ctx context.Context, path string, query url.Values, queueName string) *asynctask.AsyncTask {
	a := httptest.NewAssert(runner.t)
	ctx, err := appengine.Namespace(ctx, runner.service.Namespace())
	if err != nil {
		panic(err)
	}
	triggerPath := path
	if query != nil {
		triggerPath = fmt.Sprintf("%s?%s", triggerPath, query.Encode())
	}
	recorder := gaetest.NewRecorder(runner.service)

	// trigger the task
	var triggered asynctask.TriggerResponse
	res := recorder.TestPost(triggerPath, nil)
	a.Status(response.HTTPStatusCreated, res)
	a.JSON(&triggered, res)

	// Reqeust to the execution endpoint manually here since no module is loaded
	// on test server and the queue is not consumed automatically.
	basePath := fmt.Sprintf("%s%s.json", path, triggered.ID)
	execPath := basePath
	if query != nil {
		execPath = fmt.Sprintf("%s?%s", basePath, query.Encode())
	}
	// loop until next parameter is brank
	for {
		req := recorder.NewRequest("POST", execPath, nil)
		req.Header.Set("X-AppEngine-TaskName", queueName)
		res := recorder.TestRequest(req)
		a.Status(response.HTTPStatusOK, res)
		var next url.Values
		if strings.TrimSpace(res.Body.String()) == "true" {
			break
		}
		a.JSON(&next, res, "unexpected response from task execution endpoint %s", execPath)
		execPath = fmt.Sprintf("%s?%s", basePath, next.Encode())
	}
	// now the task has been completed.
	task := asynctask.DefaultAsyncTaskKind.MustGet(ctx, triggered.ID)
	return task
}
