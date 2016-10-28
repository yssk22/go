package gaetest

import (
	"github.com/speedland/go/web"
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
	return f()
}

func NewContext() context.Context {
	req, err := Instance().NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}
	return appengine.NewContext(req)
}

func Instance() aetest.Instance {
	return instance
}

// GET issue a GET request to the test server
func GET(path string) *web.Request {
	r, err := instance.NewRequest("GET", path, nil)
	if err != nil {
		panic(err)
	}
	return web.NewRequest(r, nil)
}
