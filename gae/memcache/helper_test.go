package memcache

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"google.golang.org/appengine"

	"github.com/speedland/go/lazy"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/x/xtesting/assert"
	"github.com/speedland/go/x/xtime"
)

func Test_CachedObjectWithExpiration(t *testing.T) {
	a := assert.New(t)
	var item, cachedItem Item
	err := CachedObjectWithExpiration(
		gaetest.NewContext(),
		"foo",
		2*time.Second,
		&item,
		func() (interface{}, error) {
			return &Item{ID: "10"}, nil
		},
		false,
	)
	a.Nil(err)
	a.EqStr("10", item.ID)
	a.Nil(Get(gaetest.NewContext(), "foo", &cachedItem))
	a.EqStr("10", cachedItem.ID)
	time.Sleep(2 * time.Second)
	a.NotNil(Get(gaetest.NewContext(), "foo", &cachedItem))
}

func Test_CacheResponseWithExpire(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	router := web.NewRouter(&web.Option{
		HMACKey: web.DefaultOption.HMACKey,
		InitContext: func(r *http.Request) context.Context {
			return appengine.NewContext(r)
		},
	})
	expires := 5 * time.Second
	router.Get("/", CacheResponseWithExpire(lazy.New("myname"), 5*time.Second, web.HandlerFunc(func(req *web.Request, _ web.NextHandler) *response.Response {
		now := xtime.Now()
		return response.NewText(fmt.Sprintf("HelloWorld - %s", now))
	})))
	recorder := gaetest.NewRecorder(router)
	resp := recorder.TestGet("/")
	a.EqInt(200, resp.Code)
	cachedBody := resp.Body.String()

	resp = recorder.TestGet("/")
	a.EqInt(200, resp.Code)
	a.EqStr(cachedBody, resp.Body.String())
	time.Sleep(expires)

	resp = recorder.TestGet("/")
	a.EqInt(200, resp.Code)
	a.Not(cachedBody == resp.Body.String())
}
