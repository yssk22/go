package gaetest

import (
	"os"
	"sync"

	"github.com/speedland/go/web"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

var instance aetest.Instance

var once sync.Once

// Exit exits the process with shutting down test server.
func Exit(code int) {
	if instance != nil {
		instance.Close()
	}
	os.Exit(code)
}

func NewContext() context.Context {
	req, err := Instance().NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}
	return appengine.NewContext(req)
}

func Instance() aetest.Instance {
	var err error
	once.Do(func() {
		instance, err = aetest.NewInstance(&aetest.Options{
			AppID: "gaetest",
			StronglyConsistentDatastore: true,
		})
		if err != nil {
			panic(err)
		}
	})
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
