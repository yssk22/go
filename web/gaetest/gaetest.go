package gaeserver

import (
	"sync"

	"github.com/speedland/go/web"

	"google.golang.org/appengine/aetest"
)

var instance aetest.Instance

var once sync.Once

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
	return web.NewRequest(r)
}
