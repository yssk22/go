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
	"google.golang.org/appengine"
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
			a.OK(e.BeforeSaveProcessed)

			_, ents, err := k.GetMulti(gaetest.NewContext(), keys)
			a.Nil(err)
			a.EqInt(1, len(keys))
			a.EqInt(1, len(ents))
			a.NotNil(ents[0])
			a.EqStr(e.ID, ents[0].ID)
			a.EqStr(e.Desc, ents[0].Desc)
			a.OK(ents[0].BeforeSaveProcessed)
		},
	)
}

func TestExampleKind_ReplaceMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_ReplaceMulti.json", nil))
	k := &ExampleKind{}
	r := ExampleKindReplacerFunc(func(e1 *Example, e2 *Example) *Example {
		if e2.Desc != "" {
			e1.Desc = e2.Desc
		}
		return e1
	})
	_, ents, err := k.ReplaceMulti(gaetest.NewContext(), []*Example{
		&Example{
			ID:   "example-1",
			Desc: "",
		},
	}, r)
	a.Nil(err)
	a.EqStr("example-1 description", ents[0].Desc)
	_, ents, err = k.ReplaceMulti(gaetest.NewContext(), []*Example{
		&Example{
			ID:   "example-1",
			Desc: "replaced",
		},
	}, r)
	a.Nil(err)
	a.EqStr("replaced", ents[0].Desc)

	_, ents, err = k.ReplaceMulti(gaetest.NewContext(), []*Example{
		&Example{
			ID:   "example-3",
			Desc: "newone",
		},
	}, r)
	a.Nil(err)
	a.EqStr("newone", ents[0].Desc)
	a.EqInt(2, NewExampleQuery().MustCount(gaetest.NewContext()))
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
	keys, values, err := q.GetAll(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(1, len(keys))
	a.EqStr("example-2", values[0].ID)
}

func TestExampleQuery_GetAll_ViaKeys(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExampleQuery_GetAll.json", nil))

	q := NewExampleQuery().ViaKeys(DefaultExampleKind)
	q.Eq("ID", lazy.New("example-2"))
	keys, values, err := q.GetAll(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(1, len(keys))
	a.EqStr("example-2", values[0].ID)

	var e1 Example
	err = memcache.Get(gaetest.NewContext(), ent.GetMemcacheKey(keys[0]), &e1)
	a.Nil(err)
	a.EqStr(values[0].Desc, e1.Desc)
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

func TestExampleQuery_Run_ViaKeys(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExampleQuery_Run.json", nil))

	var e Example
	q := NewExampleQuery().Asc("ID").Limit(lazy.New(2)).ViaKeys(DefaultExampleKind)
	p, err := q.Run(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(2, len(p.Data))
	a.EqStr("example-1", p.Data[0].ID)
	a.EqStr("example-2", p.Data[1].ID)
	a.EqStr("", p.Start)
	// check cache
	a.Nil(memcache.Get(gaetest.NewContext(), ent.GetMemcacheKey(p.Keys[0]), &e))
	a.Nil(memcache.Get(gaetest.NewContext(), ent.GetMemcacheKey(p.Keys[1]), &e))
	next := p.End

	q = NewExampleQuery().Asc("ID").Limit(lazy.New(2)).Start(lazy.New(next)).ViaKeys(DefaultExampleKind)
	p, err = q.Run(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(2, len(p.Data))
	a.EqStr("example-3", p.Data[0].ID)
	a.EqStr("example-4", p.Data[1].ID)
	a.Nil(memcache.Get(gaetest.NewContext(), ent.GetMemcacheKey(p.Keys[0]), &e))
	a.Nil(memcache.Get(gaetest.NewContext(), ent.GetMemcacheKey(p.Keys[1]), &e))
	a.EqStr(next, p.Start)

	q = NewExampleQuery().Asc("ID").Limit(lazy.New(2)).Start(lazy.New(p.Start)).ViaKeys(DefaultExampleKind)
	p, err = q.Run(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(2, len(p.Data))
	a.EqStr("example-3", p.Data[0].ID)
	a.EqStr("example-4", p.Data[1].ID)
	a.Nil(memcache.Get(gaetest.NewContext(), ent.GetMemcacheKey(p.Keys[0]), &e))
	a.Nil(memcache.Get(gaetest.NewContext(), ent.GetMemcacheKey(p.Keys[1]), &e))

	q = NewExampleQuery().Asc("ID").Limit(lazy.New(2)).Start(lazy.New(p.End)).ViaKeys(DefaultExampleKind)
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

func TestExampleKind_SearchKeys(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_Search.json", nil))
	examples, err := NewExampleQuery().Asc("ID").GetAllValues(gaetest.NewContext())
	a.Nil(err)
	_, err = DefaultExampleKind.PutMulti(gaetest.NewContext(), examples)
	a.Nil(err)

	p, err := DefaultExampleKind.SearchKeys(gaetest.NewContext(), "Desc: example-2", nil)
	a.Nil(err)
	a.EqInt(1, len(p.Keys))
}

func TestExampleKind_SearchValues(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_Search.json", nil))
	examples, err := NewExampleQuery().Asc("ID").GetAllValues(gaetest.NewContext())
	a.Nil(err)
	_, err = DefaultExampleKind.PutMulti(gaetest.NewContext(), examples)
	a.Nil(err)

	p, err := DefaultExampleKind.SearchValues(gaetest.NewContext(), "Desc: example-2", nil)
	a.Nil(err)
	a.EqInt(1, len(p.Data))
	a.EqStr("example-2 description", p.Data[0].Desc)
}

func TestExampleKind_DeleteMatched(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_DeleteMatched.json", nil))
	a.EqInt(4, NewExampleQuery().MustCount(gaetest.NewContext()))
	q := NewExampleQuery().Le("Digit", lazy.New(2))
	a.EqInt(2, q.MustCount(gaetest.NewContext()))
	deleted, err := DefaultExampleKind.DeleteMatched(gaetest.NewContext(), q)
	a.Nil(err)
	a.EqInt(2, deleted)
	a.EqInt(2, NewExampleQuery().MustCount(gaetest.NewContext()))
	a.Nil(DefaultExampleKind.MustGet(gaetest.NewContext(), "example-1"))
	a.Nil(DefaultExampleKind.MustGet(gaetest.NewContext(), "example-2"))
	a.NotNil(DefaultExampleKind.MustGet(gaetest.NewContext(), "example-3"))
	a.NotNil(DefaultExampleKind.MustGet(gaetest.NewContext(), "example-4"))
}

func TestExampleKind_EnforceNamespace(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.CleanupStorage(gaetest.NewContext(), "", "myns"))
	ctx, err := appengine.Namespace(gaetest.NewContext(), "myns")
	a.Nil(err)

	kind := &ExampleKind{}
	kind.EnforceNamespace("myns", true)

	// Get
	a.NotNil(kind.MustPut(gaetest.NewContext(), &Example{
		ID: "myid",
	}))
	ent := DefaultExampleKind.MustGet(ctx, "myid")
	a.NotNil(ent)
	ent = kind.MustGet(gaetest.NewContext(), "myid")
	a.NotNil(ent)

	// Put (Update)
	a.NotNil(kind.MustPut(gaetest.NewContext(), &Example{
		ID:   "myid",
		Desc: "my description",
	}))
	ent = DefaultExampleKind.MustGet(ctx, "myid")
	a.NotNil(ent)
	a.EqStr("my description", ent.Desc)
	ent = kind.MustGet(ctx, "myid")
	a.NotNil(ent)
	a.EqStr("my description", ent.Desc)

	// Delete
	a.NotNil(kind.MustDelete(gaetest.NewContext(), "myid"))
	ent = DefaultExampleKind.MustGet(ctx, "myid")
	a.Nil(ent)
	ent = kind.MustGet(ctx, "myid")
	a.Nil(ent)
}
