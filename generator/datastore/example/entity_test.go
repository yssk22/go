package example

import (
	"os"
	"testing"
	"time"

	ds "github.com/yssk22/go/gae/datastore"
	"github.com/yssk22/go/gae/gaetest"
	"github.com/yssk22/go/gae/memcache"
	"github.com/yssk22/go/x/xtesting/assert"
	"github.com/yssk22/go/x/xtime"
	memc "google.golang.org/appengine/memcache"
)

var entityStore = GetEntityDatastore()

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestEntity_Get(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_Get.json", nil))

	_, value, err := entityStore.Get(gaetest.NewContext(), "entity-1")
	a.Nil(err)
	a.NotNil(value)
	a.EqStr("entity-1 description", value.Desc)
}

func TestEntity_GetMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_GetMulti.json", nil))

	keys, values, err := entityStore.GetMulti(gaetest.NewContext(), []string{"entity-1", "entity-2"})
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

	keys, values, err := entityStore.GetMulti(gaetest.NewContext(), []string{"aaa", "entity-2"})
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

	keys, values, err := entityStore.GetMulti(gaetest.NewContext(), []string{"entity-1", "not-exists"})
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
	a.OK(memc.ErrCacheMiss == err)

	// Delete datastore (to check cache can work)
	a.Nil(gaetest.CleanupDatastore(gaetest.NewContext()))
	_, values, err = entityStore.GetMulti(gaetest.NewContext(), []string{"entity-1"})
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
			keys, err := entityStore.PutMulti(gaetest.NewContext(), []*Entity{e})
			a.Nil(err)
			a.EqInt(1, len(keys))
			a.EqStr(e.ID, keys[0].StringID())
			a.EqTime(now, e.UpdatedAt)

			_, ents, err := entityStore.GetMulti(gaetest.NewContext(), keys)
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
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_DeleteMulti.json", nil))

	keys, err := entityStore.DeleteMulti(gaetest.NewContext(), []string{"entity-1", "entity-2"})
	a.Nil(err)
	a.EqInt(2, len(keys))
	_, ents := entityStore.MustGetMulti(gaetest.NewContext(), []string{"entity-1", "entity-2"})
	a.Nil(err)
	a.Nil(ents[0])
	a.Nil(ents[1])
}

func TestExampleKind_ReplaceMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.ResetMemcache(gaetest.NewContext()))
	a.Nil(gaetest.ResetFixtureFromFile(gaetest.NewContext(), "./fixture/TestEntity_ReplaceMulti.json", nil))
	r := EntityReplacerFunc(func(e1 *Entity, e2 *Entity) *Entity {
		if e2.Desc != "" {
			e1.Desc = e2.Desc
		}
		return e1
	})
	_, ents, err := entityStore.ReplaceMulti(gaetest.NewContext(), []*Entity{
		&Entity{
			ID:   "entity-1",
			Desc: "",
		},
	}, r)
	a.Nil(err)
	a.EqStr("entity-1 description", ents[0].Desc)
	_, ents, err = entityStore.ReplaceMulti(gaetest.NewContext(), []*Entity{
		&Entity{
			ID:   "entity-1",
			Desc: "replaced",
		},
	}, r)
	a.Nil(err)
	a.EqStr("replaced", ents[0].Desc)

	_, ents, err = entityStore.ReplaceMulti(gaetest.NewContext(), []*Entity{
		&Entity{
			ID:   "entity-3",
			Desc: "newone",
		},
	}, r)
	a.Nil(err)
	a.EqStr("newone", ents[0].Desc)
	// a.EqInt(2, NewExampleQuery().MustCount(gaetest.NewContext()))
}
