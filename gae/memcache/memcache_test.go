package memcache

import (
	"os"
	"testing"

	"github.com/yssk22/go/gae/gaetest"
	"github.com/yssk22/go/x/xtesting/assert"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

type Item struct {
	ID string
}

func TestSet(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))

	item := &Item{ID: "FOO"}
	a.Nil(Set(gaetest.NewContext(), "a", item))
	a.OK(Exists(gaetest.NewContext(), "a"))
}

func TestGet(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))

	item := &Item{ID: "FOO"}
	a.Nil(Set(gaetest.NewContext(), "a", item))

	item2 := &Item{}
	a.Nil(Get(gaetest.NewContext(), "a", item2))
	a.EqStr("FOO", item2.ID)
}

func TestGet_cacheMiss(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))

	item2 := &Item{}
	err := Get(gaetest.NewContext(), "a", item2)
	a.OK(memcache.ErrCacheMiss == err)
}

func TestGetMulti_cacheMiss(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))

	item := &Item{ID: "FOO"}
	a.Nil(Set(gaetest.NewContext(), "a", item))

	items := make([]*Item, 2, 2)
	err := GetMulti(gaetest.NewContext(), []string{
		"a", "b",
	}, items)
	merr, ok := err.(appengine.MultiError)
	a.OK(ok)
	a.Nil(merr[0])
	a.OK(memcache.ErrCacheMiss == merr[1])
	a.EqStr("FOO", items[0].ID)
	a.Nil(items[1])
}
