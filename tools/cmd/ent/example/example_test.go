package example

import (
	"os"
	"testing"
	"time"

	"github.com/speedland/go/x/xtime"

	"github.com/speedland/go/gae/datastore/ent"
	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/gae/memcache"

	"github.com/speedland/go/x/xtesting/assert"
	memc "google.golang.org/appengine/memcache"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestEample_NewKey(t *testing.T) {
	a := assert.New(t)
	n := (&ExampleKind{}).New()
	n.ID = "FOO"
	key := n.NewKey(gaetest.NewContext())
	a.EqStr("Example", key.Kind())
	a.EqStr("FOO", key.StringID())
}

func TestExampleKind_New(t *testing.T) {
	a := assert.New(t)
	now := time.Date(2016, 1, 1, 0, 0, 0, 0, xtime.JST)
	xtime.RunAt(
		now,
		func() {
			n := (&ExampleKind{}).New()
			a.EqStr("This is default value", n.Desc)
			a.EqInt(10, n.Digit)
			a.EqTime(now, n.CreatedAt)
			a.EqTime(
				time.Date(2016, 01, 01, 20, 12, 10, 0, time.UTC), //2016-01-01T20:12:10Z
				n.DefaultTime,
			)
		},
	)
}

func TestExampleKind_Get(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_Get.json", nil))

	k := &ExampleKind{}
	_, value, err := k.Get(gaetest.NewContext(), "example-1")
	a.Nil(err)
	a.NotNil(value)
	a.EqStr("example-1 description", value.Desc)
}

func TestExampleKind_GetMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_GetMulti.json", nil))

	k := &ExampleKind{}
	keys, values, err := k.GetMulti(gaetest.NewContext(), "example-1", "example-2")
	a.Nil(err)
	a.EqInt(2, len(keys))
	a.EqInt(2, len(values))
	a.NotNil(values[0])
	a.NotNil(values[1])
	a.EqStr("example-1 description", values[0].Desc)
	a.EqStr("example-2 description", values[1].Desc)
}

func TestExampleKind_GetMulti_notFound(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_GetMulti.json", nil))

	k := &ExampleKind{}
	keys, values, err := k.GetMulti(gaetest.NewContext(), "aaa", "example-2")
	a.Nil(err)
	a.EqInt(2, len(keys))
	a.EqInt(2, len(values))
	a.Nil(values[0])
	a.NotNil(values[1])
}

func TestExampleKind_GetMulti_useDefaultIfNil(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_GetMulti.json", nil))

	k := (&ExampleKind{}).UseDefaultIfNil(true)
	keys, values, err := k.GetMulti(gaetest.NewContext(), "aaa", "example-2")
	a.Nil(err)
	a.EqInt(2, len(keys))
	a.EqInt(2, len(values))
	a.NotNil(values[0])
	a.EqStr("aaa", values[0].ID)
	a.EqStr("This is default value", values[0].Desc)
	a.NotNil(values[1])
}

func TestExampleKind_GetMulti_cacheCreation(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_GetMulti.json", nil))

	k := &ExampleKind{}
	keys, values, err := k.GetMulti(gaetest.NewContext(), "example-1", "not-exists")
	a.Nil(err)
	a.EqInt(2, len(keys))
	a.EqInt(2, len(values))
	a.NotNil(values[0])
	a.Nil(values[1])

	// Check cache
	var e1 Example
	err = memcache.Get(gaetest.NewContext(), ent.GetMemcacheKey(keys[0]), &e1)
	a.Nil(err)
	a.EqStr(values[0].Desc, e1.Desc)

	err = memcache.Get(gaetest.NewContext(), ent.GetMemcacheKey(keys[1]), &e1)
	a.OK(memc.ErrCacheMiss == err)

	// Delete datastore (to check cache can work)
	a.Nil(gaetest.CleanupDatastore(gaetest.NewContext()))
	_, values, _ = k.GetMulti(gaetest.NewContext(), "example-1")
	a.NotNil(values[0])
	a.EqStr(e1.Desc, values[0].Desc)
}
