package gaeserver

import (
	"sync"

	"github.com/speedland/go/web"

	"google.golang.org/appengine/aetest"
)

var instance aetest.Instance

func Instance() aetest.Instance {
	sync.Once(func() {
		instance = aetest.NewInstance(&aetest.Options{
			AppID: "gaetest",
			StronglyConsistentDatastore: true,
		})
	})
	return instance
}

// GET issue a GET request to the test server
func GET(path string) *web.Request {
	r := instance.NewRequest("GET", path, nil)
}
