package gaetest

import (
	"io"
	"net/http"

	"github.com/speedland/go/web/httptest"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

var instance aetest.Instance

// Run runs gaetest server and run f in a clean way.
// Note:
//    this does not kill dev_appserver process if the test code fails in panic
//    even we use defer to close the instance. This is because each test are executed
//    in the different goroutine, which cannot be recovered.
func Run(f func() int) int {
	var err error
	defer func() {
		if instance != nil {
			instance.Close()
		}
	}()
	instance, err = aetest.NewInstance(&aetest.Options{
		AppID: "gaetest",
		StronglyConsistentDatastore: true,
	})
	if err != nil {
		panic(err)
	}
	prepareTaskQueueInTest()
	return f()
}

// NewContext returns a new appengine context.Context
func NewContext() context.Context {
	req, err := Instance().NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}
	return appengine.NewContext(req)
}

// Instance returns a current appengine instance
func Instance() aetest.Instance {
	if instance == nil {
		panic("GAE test instance is not initialized. Call gaetest.Run() on your TestMain function.")
	}
	return instance
}

// NewRequest returns a new *http.Request bound with appengine context.Context
func NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	return Instance().NewRequest(method, path, body)
}

// NewRecorder returns a new *httptest.Recorder object
func NewRecorder(handler http.Handler) *httptest.Recorder {
	return httptest.NewRecorderWithFactory(
		handler,
		httptest.RequestFactoryFunc(NewRequest),
	)
}
