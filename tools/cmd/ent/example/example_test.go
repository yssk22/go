package example

import (
	"encoding/json"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/lazy"
	"github.com/speedland/go/x/xtime"

	"github.com/speedland/go/ent"
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
	n := NewExample()
	n.ID = "FOO"
	key := n.NewKey(gaetest.NewContext())
	a.EqStr("Example", key.Kind())
	a.EqStr("FOO", key.StringID())
}

func TestExample_UpdateByForm(t *testing.T) {
	a := assert.New(t)
	n := NewExample()
	n.Desc = "foo"
	form := url.Values{
		"desc":        []string{"val"},
		"custom_type": []string{"#ff0000"},
	}
	getter := keyvalue.GetterFunc(func(key interface{}) (interface{}, error) {
		v, ok := form[key.(string)]
		if ok && len(v) >= 0 {
			return v[0], nil
		}
		return nil, keyvalue.KeyError(key.(string))
	})
	n.UpdateByForm(keyvalue.NewGetProxy(getter))
	a.EqStr("val", n.Desc)
	a.EqStr("#ff0000", n.CustomType.ToHexString())
}

func TestExampleKind_New(t *testing.T) {
	a := assert.New(t)
	now := time.Date(2016, 1, 1, 0, 0, 0, 0, xtime.JST)
	xtime.RunAt(
		now,
		func() {
			n := NewExample()
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
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_Get.json", nil))

	k := &ExampleKind{}
	_, value, err := k.Get(gaetest.NewContext(), "example-1")
	a.Nil(err)
	a.NotNil(value)
	a.EqStr("example-1 description", value.Desc)
}

func TestExampleKind_GetMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_GetMulti.json", nil))

	k := &ExampleKind{}
	keys, values, err := k.GetMulti(gaetest.NewContext(), []string{"example-1", "example-2"})
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
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_GetMulti.json", nil))

	k := &ExampleKind{}
	keys, values, err := k.GetMulti(gaetest.NewContext(), []string{"aaa", "example-2"})
	a.Nil(err)
	a.EqInt(2, len(keys))
	a.EqInt(2, len(values))
	a.Nil(values[0])
	a.NotNil(values[1])
}

func TestExampleKind_GetMulti_useDefaultIfNil(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_GetMulti.json", nil))

	k := (&ExampleKind{}).UseDefaultIfNil(true)
	keys, values, err := k.GetMulti(gaetest.NewContext(), []string{"aaa", "example-2"})
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
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_GetMulti.json", nil))

	k := &ExampleKind{}
	keys, values, err := k.GetMulti(gaetest.NewContext(), []string{"example-1", "not-exists"})
	a.Nil(err)
	a.EqInt(2, len(keys), "%v, %v", keys, values)
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
	_, values, err = k.GetMulti(gaetest.NewContext(), []string{"example-1"})
	a.Nil(err)
	a.NotNil(values[0])
	a.EqStr(e1.Desc, values[0].Desc)
}

func TestExampleKind_PutMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))

	k := &ExampleKind{}
	e := NewExample()
	e.ID = "foo"
	e.Desc = "PUT TEST"

	now := time.Date(2016, 1, 1, 12, 12, 0, 0, xtime.JST)
	xtime.RunAt(
		now,
		func() {
			keys, err := k.PutMulti(gaetest.NewContext(), []*Example{e})
			a.Nil(err)
			a.EqInt(1, len(keys))
			a.EqStr(e.ID, keys[0].StringID())
			a.EqTime(now, e.UpdatedAt)

			_, ents, err := k.GetMulti(gaetest.NewContext(), keys)
			a.Nil(err)
			a.EqInt(1, len(keys))
			a.EqInt(1, len(ents))
			a.NotNil(ents[0])
			a.EqStr(e.ID, ents[0].ID)
			a.EqStr(e.Desc, ents[0].Desc)
		},
	)
}

func TestExampleKind_DeleteMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_DeleteMulti.json", nil))

	k := &ExampleKind{}
	keys, err := k.DeleteMulti(gaetest.NewContext(), []string{"example-1", "example-2"})
	a.Nil(err)
	a.EqInt(2, len(keys))
	ents := k.MustGetMulti(gaetest.NewContext(), []string{"example-1", "example-2"})
	a.Nil(err)
	a.Nil(ents[0])
	a.Nil(ents[1])
}

func TestExampleQuery_GetAll(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExampleQuery_GetAll.json", nil))

	q := NewExampleQuery()
	q.Eq("ID", lazy.New("example-2"))
	keys, value, err := q.GetAll(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(1, len(keys))
	a.EqStr("example-2", value[0].ID)
}

func TestExampleQuery_Run(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExampleQuery_Run.json", nil))

	q := NewExampleQuery().Asc("ID").Limit(lazy.New(2))
	p, err := q.Run(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(2, len(p.Data))
	a.EqStr("example-1", p.Data[0].ID)
	a.EqStr("example-2", p.Data[1].ID)
	a.EqStr("", p.Start)
	next := p.End

	q = NewExampleQuery().Asc("ID").Limit(lazy.New(2)).Start(lazy.New(next))
	p, err = q.Run(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(2, len(p.Data))
	a.EqStr("example-3", p.Data[0].ID)
	a.EqStr("example-4", p.Data[1].ID)
	a.EqStr(next, p.Start)

	q = NewExampleQuery().Asc("ID").Limit(lazy.New(2)).Start(lazy.New(p.Start))
	p, err = q.Run(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(2, len(p.Data))
	a.EqStr("example-3", p.Data[0].ID)
	a.EqStr("example-4", p.Data[1].ID)

	q = NewExampleQuery().Asc("ID").Limit(lazy.New(2)).Start(lazy.New(p.End))
	p, err = q.Run(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(0, len(p.Data))
}

func TestExamplePagination_MarshalJSON(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExampleQuery_Run.json", nil))

	q := NewExampleQuery().Limit(lazy.New(0))
	p, err := q.Run(gaetest.NewContext())
	a.Nil(err)
	b, _ := json.Marshal(p)
	a.EqByteString(`{"start":"","end":"","data":[]}`, b)
}
