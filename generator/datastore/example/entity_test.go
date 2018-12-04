package example

import (
	"log"
	"os"
	"testing"
	"time"

	ds "github.com/yssk22/go/gae/datastore"
	"github.com/yssk22/go/gae/gaetest"
	"github.com/yssk22/go/gae/memcache"
	"github.com/yssk22/go/x/xtesting/assert"
	"github.com/yssk22/go/x/xtime"
)


func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestEntity_Get(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_Get.json", nil))

	_, value, err := entityKindInstance.Get(gaetest.NewContext(), "entity-1")
	a.Nil(err)
	a.NotNil(value)
	a.EqStr("entity-1 description", value.Desc)
}

func TestEntity_GetMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_GetMulti.json", nil))

	keys, values, err := entityKindInstance.GetMulti(gaetest.NewContext(), []string{"entity-1", "entity-2"})
	a.Nil(err)
	a.EqInt(2, len(keys))
	a.EqInt(2, len(values))
	a.NotNil(values[0])
	a.NotNil(values[1])
	a.EqStr("entity-1 description", values[0].Desc)
	a.EqStr("entity-2 description", values[1].Desc)
}

func TestEntity_GetMulti_notFound(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_GetMulti.json", nil))

	keys, values, err := entityKindInstance.GetMulti(gaetest.NewContext(), []string{"aaa", "entity-2"})
	a.Nil(err)
	a.EqInt(2, len(keys))
	a.EqInt(2, len(values))
	a.Nil(values[0])
	a.NotNil(values[1])
}

func TestEntity_GetMulti_cacheCreation(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_GetMulti.json", nil))

	keys, values, err := entityKindInstance.GetMulti(gaetest.NewContext(), []string{"entity-1", "not-exists"})
	a.Nil(err)
	a.EqInt(2, len(keys), "%v, %v", keys, values)
	a.EqInt(2, len(values))
	a.NotNil(values[0])
	a.Nil(values[1])

	// Check cache
	var e1 Entity
	err = memcache.Get(gaetest.NewContext(), ds.GetMemcacheKey(keys[0]), &e1)
	a.Nil(err)
	a.EqStr(values[0].Desc, e1.Desc)

	err = memcache.Get(gaetest.NewContext(), ds.GetMemcacheKey(keys[1]), &e1)
	a.Nil(err)

	// Delete datastore (to check cache can work)
	a.Nil(gaetest.CleanupDatastore(gaetest.NewContext()))
	_, values, err = entityKindInstance.GetMulti(gaetest.NewContext(), []string{"entity-1"})
	a.Nil(err)
	a.NotNil(values[0])
	a.EqStr(e1.Desc, values[0].Desc)
}

func TestEntity_PutMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))

	e := &Entity{}
	e.ID = "foo"
	e.Desc = "PUT TEST"

	now := time.Date(2016, 1, 1, 12, 12, 0, 0, xtime.JST)
	xtime.RunAt(
		now,
		func() {
			keys, err := entityKindInstance.PutMulti(gaetest.NewContext(), []*Entity{e})
			a.Nil(err)
			a.EqInt(1, len(keys))
			a.EqStr(e.ID, keys[0].StringID())
			a.EqTime(now, e.UpdatedAt)

			_, ents, err := entityKindInstance.GetMulti(gaetest.NewContext(), keys)
			a.Nil(err)
			a.EqInt(1, len(keys))
			a.EqInt(1, len(ents))
			a.NotNil(ents[0])
			a.EqStr(e.ID, ents[0].ID)
			a.EqStr(e.Desc, ents[0].Desc)
		},
	)
}

func TestEntity_DeleteMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_DeleteMulti.json", nil))

	keys, err := entityKindInstance.DeleteMulti(gaetest.NewContext(), []string{"entity-1", "entity-2"})
	a.Nil(err)
	a.EqInt(2, len(keys))
	_, ents := entityKindInstance.MustGetMulti(gaetest.NewContext(), []string{"entity-1", "entity-2"})
	a.Nil(err)
	a.Nil(ents[0])
	a.Nil(ents[1])
}

func TestEntity_ReplaceMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_ReplaceMulti.json", nil))
	r := EntityReplacerFunc(func(e1 *Entity, e2 *Entity) *Entity {
		if e2.Desc != "" {
			e1.Desc = e2.Desc
		}
		return e1
	})
	_, ents, err := entityKindInstance.ReplaceMulti(gaetest.NewContext(), []*Entity{
		&Entity{
			ID:   "entity-1",
			Desc: "",
		},
	}, r)
	a.Nil(err)
	a.EqStr("entity-1 description", ents[0].Desc)
	_, ents, err = entityKindInstance.ReplaceMulti(gaetest.NewContext(), []*Entity{
		&Entity{
			ID:   "entity-1",
			Desc: "replaced",
		},
	}, r)
	a.Nil(err)
	a.EqStr("replaced", ents[0].Desc)

	_, ents, err = entityKindInstance.ReplaceMulti(gaetest.NewContext(), []*Entity{
		&Entity{
			ID:   "entity-3",
			Desc: "newone",
		},
	}, r)
	a.Nil(err)
	a.EqStr("newone", ents[0].Desc)
	a.EqInt(2, NewEntityQuery().MustCount(gaetest.NewContext()))
}

func TestEntity_GetAll(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_GetAll.json", nil))
	
	keys, values, err := NewEntityQuery().EqID("entity-2").GetAll(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(1, len(keys))
	a.EqStr("entity-2", values[0].ID)
	a.EqStr("entity-2 description", values[0].Desc)
}

func TestEntity_GetAll_ViaKeys(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_GetAll.json", nil))

	keys, values, err := NewEntityQuery().EqID("entity-2").ViaKeys().GetAll(gaetest.NewContext())
	a.Nil(err)
	a.EqInt(1, len(keys))
	a.EqStr("entity-2", values[0].ID)
	a.EqStr("entity-2 description", values[0].Desc)

	var e1 Entity
	err = memcache.Get(gaetest.NewContext(), ds.GetMemcacheKey(keys[0]), &e1)
	a.Nil(err)
	a.EqStr(values[0].ID, e1.ID)
	a.EqStr(values[0].Desc, e1.Desc)
}

func TestEntity_Run(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_GetAll.json", nil))

	iter, err := NewEntityQuery().EqID("entity-2").Run(gaetest.NewContext())
	a.Nil(err)
	_, value := iter.MustNext()
	a.Nil(err)
	a.EqStr("entity-2", value.ID)
	a.EqStr("entity-2 description", value.Desc)
	_, value = iter.MustNext()
	a.Nil(value)
}

func TestEntity_Run_ViaKeys(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_GetAll.json", nil))

	iter, err := NewEntityQuery().EqID("entity-2").ViaKeys().Run(gaetest.NewContext())
	a.Nil(err)
	key, value := iter.MustNext()
	a.Nil(err)
	a.EqStr("entity-2", value.ID)
	a.EqStr("entity-2 description", value.Desc)

	var e1 Entity
	err = memcache.Get(gaetest.NewContext(), ds.GetMemcacheKey(key), &e1)
	a.Nil(err)
	a.EqStr(value.ID, e1.ID)
	a.EqStr(value.Desc, e1.Desc)
}

func TestEntity_RunAll(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_GetAll.json", nil))

	_, values, next := NewEntityQuery().EqID("entity-2").ViaKeys().MustRunAll(gaetest.NewContext())
	a.EqInt(1, len(values))
	a.EqStr("entity-2", values[0].ID)
	a.EqStr("entity-2 description", values[0].Desc)
	a.EqStr("CikSI2oLZGV2fmdhZXRlc3RyFAsSBkVudGl0eSIIZW50aXR5LTIMGAAgAA", next)
}

func TestEntity_DeleteMatched(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_DeleteMatched.json", nil))
	a.EqInt(4, NewEntityQuery().MustCount(gaetest.NewContext()))
	q := NewEntityQuery().LeDigit(2)
	a.EqInt(2, q.MustCount(gaetest.NewContext()))
	deleted := entityKindInstance.MustDeleteMatched(gaetest.NewContext(), q)
	a.EqInt(2, len(deleted))
	a.EqInt(2, NewEntityQuery().MustCount(gaetest.NewContext()))
	log.Println(NewEntityQuery().MustGetAll(gaetest.NewContext()))
	_, value := entityKindInstance.MustGet(gaetest.NewContext(), "entity-1")
	a.Nil(value)
	_, value = entityKindInstance.MustGet(gaetest.NewContext(), "entity-2")
	a.Nil(value)
	_, value = entityKindInstance.MustGet(gaetest.NewContext(), "entity-3")
	a.NotNil(value)
	_, value = entityKindInstance.MustGet(gaetest.NewContext(), "entity-4")
	a.NotNil(value)
}


// func TestExampleKind_SearchKeys(t *testing.T) {
// 	a := assert.New(t)
// 	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
// 	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_Search.json", nil))
// 	examples, err := NewExampleQuery().Asc("ID").GetAllValues(gaetest.NewContext())
// 	a.Nil(err)
// 	_, err = DefaultExampleKind.PutMulti(gaetest.NewContext(), examples)
// 	a.Nil(err)

// 	p, err := DefaultExampleKind.SearchKeys(gaetest.NewContext(), "Desc: example-2", nil)
// 	a.Nil(err)
// 	a.EqInt(1, len(p.Keys))
// }

// func TestExampleKind_SearchValues(t *testing.T) {
// 	a := assert.New(t)
// 	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
// 	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_Search.json", nil))
// 	examples, err := NewExampleQuery().Asc("ID").GetAllValues(gaetest.NewContext())
// 	a.Nil(err)
// 	_, err = DefaultExampleKind.PutMulti(gaetest.NewContext(), examples)
// 	a.Nil(err)

// 	p, err := DefaultExampleKind.SearchValues(gaetest.NewContext(), "Desc: example-2", nil)
// 	a.Nil(err)
// 	a.EqInt(1, len(p.Data))
// 	a.EqStr("example-2 description", p.Data[0].Desc)
// }

// func TestExampleKind_EnforceNamespace(t *testing.T) {
// 	a := assert.New(t)
// 	a.Nil(gaetest.CleanupStorage(gaetest.NewContext(), "", "myns"))
// 	ctx, err := appengine.Namespace(gaetest.NewContext(), "myns")
// 	a.Nil(err)

// 	kind := &ExampleKind{}
// 	kind.EnforceNamespace("myns", true)

// 	// Get
// 	a.NotNil(kind.MustPut(gaetest.NewContext(), &Example{
// 		ID: "myid",
// 	}))
// 	ent := DefaultExampleKind.MustGet(ctx, "myid")
// 	a.NotNil(ent)
// 	ent = kind.MustGet(gaetest.NewContext(), "myid")
// 	a.NotNil(ent)

// 	// Put (Update)
// 	a.NotNil(kind.MustPut(gaetest.NewContext(), &Example{
// 		ID:   "myid",
// 		Desc: "my description",
// 	}))
// 	ent = DefaultExampleKind.MustGet(ctx, "myid")
// 	a.NotNil(ent)
// 	a.EqStr("my description", ent.Desc)
// 	ent = kind.MustGet(ctx, "myid")
// 	a.NotNil(ent)
// 	a.EqStr("my description", ent.Desc)

// 	// Delete
// 	a.NotNil(kind.MustDelete(gaetest.NewContext(), "myid"))
// 	ent = DefaultExampleKind.MustGet(ctx, "myid")
// 	a.Nil(ent)
// 	ent = kind.MustGet(ctx, "myid")
// 	a.Nil(ent)
// }
