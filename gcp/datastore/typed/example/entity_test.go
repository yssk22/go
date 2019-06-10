package example

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	// ds "github.com/yssk22/go/gae/datastore"
	// "github.com/yssk22/go/gae/gaetest"
	// "github.com/yssk22/go/gae/memcache"

	"github.com/yssk22/go/gcp/datastore"
	"github.com/yssk22/go/x/xtesting/assert"
	"github.com/yssk22/go/x/xtime"
)

var testEnvironment *datastore.TestEnviornment
var testClient *EntityKindClient

func TestMain(m *testing.M) {
	testEnvironment = datastore.MustNewTestEnviornment()
	testClient = NewEntityKindClient(testEnvironment.GetClient())

	state := m.Run()
	err := testEnvironment.Shutdown()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(state)
}

func TestEntityKindClient(t *testing.T) {
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		a := assert.New(t)
		a.Nil(testEnvironment.Reset())
		a.Nil(testEnvironment.LoadFixture("./fixture/TestEntity_Get.json"))

		_, value, err := testClient.Get(ctx, "entity-1")
		a.Nil(err)
		a.NotNil(value)
		a.EqStr("entity-1 description", value.Desc)
	})

	t.Run("GetMulti", func(t *testing.T) {
		a := assert.New(t)
		a.Nil(testEnvironment.Reset())
		a.Nil(testEnvironment.LoadFixture("./fixture/TestEntity_GetMulti.json"))

		keys, values, err := testClient.GetMulti(ctx, []string{"entity-1", "entity-2"})
		a.Nil(err)
		a.EqInt(2, len(keys))
		a.EqInt(2, len(values))
		a.NotNil(values[0])
		a.NotNil(values[1])
		a.EqStr("entity-1 description", values[0].Desc)
		a.EqStr("entity-2 description", values[1].Desc)

		keys, values, err = testClient.GetMulti(ctx, []string{"aaa", "entity-2"})
		a.Nil(err)
		a.EqInt(2, len(keys))
		a.EqInt(2, len(values))
		a.Nil(values[0])
		a.NotNil(values[1])
	})

	t.Run("PutMulti", func(t *testing.T) {
		a := assert.New(t)
		a.Nil(testEnvironment.Reset())
		e := &Entity{}
		e.ID = "foo"
		e.Desc = "PUT TEST"

		now := time.Date(2016, 1, 1, 12, 12, 0, 0, xtime.JST)
		xtime.RunAt(
			now,
			func() {
				keys, err := testClient.PutMulti(ctx, []*Entity{e})
				a.Nil(err)
				a.EqInt(1, len(keys))
				a.EqStr(e.ID, keys[0].Name)
				a.EqTime(now, e.UpdatedAt)

				_, ents, err := testClient.GetMulti(ctx, keys)
				a.Nil(err)
				a.EqInt(1, len(keys))
				a.EqInt(1, len(ents))
				a.NotNil(ents[0])
				a.EqStr(e.ID, ents[0].ID)
				a.EqStr(e.Desc, ents[0].Desc)
			},
		)
	})

	t.Run("DeleteMulti", func(t *testing.T) {
		a := assert.New(t)
		a.Nil(testEnvironment.Reset())
		a.Nil(testEnvironment.LoadFixture("./fixture/TestEntity_DeleteMulti.json"))

		keys, err := testClient.DeleteMulti(ctx, []string{"entity-1", "entity-2"})
		a.Nil(err)
		a.EqInt(2, len(keys))
		_, ents := testClient.MustGetMulti(ctx, []string{"entity-1", "entity-2"})
		a.Nil(err)
		a.Nil(ents[0])
		a.Nil(ents[1])
	})

	t.Run("ReplaceMulti", func(t *testing.T) {
		a := assert.New(t)
		a.Nil(testEnvironment.Reset())
		a.Nil(testEnvironment.LoadFixture("./fixture/TestEntity_ReplaceMulti.json"))
		r := EntityReplacerFunc(func(e1 *Entity, e2 *Entity) *Entity {
			if e2.Desc != "" {
				e1.Desc = e2.Desc
			}
			return e1
		})
		_, ents, err := testClient.ReplaceMulti(ctx, []*Entity{
			&Entity{
				ID:   "entity-1",
				Desc: "",
			},
		}, r)
		a.Nil(err)
		a.EqStr("entity-1 description", ents[0].Desc)
		_, ents, err = testClient.ReplaceMulti(ctx, []*Entity{
			&Entity{
				ID:   "entity-1",
				Desc: "replaced",
			},
		}, r)
		a.Nil(err)
		a.EqStr("replaced", ents[0].Desc)

		_, ents, err = testClient.ReplaceMulti(ctx, []*Entity{
			&Entity{
				ID:   "entity-3",
				Desc: "newone",
			},
		}, r)
		a.Nil(err)
		a.EqStr("newone", ents[0].Desc)
		a.EqInt(2, testClient.MustCount(ctx, NewEntityQuery()))
	})

	t.Run("Query", func(t *testing.T) {
		t.Run("GetAll", func(t *testing.T) {
			a := assert.New(t)
			a.Nil(testEnvironment.Reset())
			a.Nil(testEnvironment.LoadFixture("./fixture/TestEntity_GetAll.json"))

			keys, values, err := testClient.GetAll(ctx, NewEntityQuery().EqID("entity-2"))
			a.Nil(err)
			a.EqInt(1, len(keys))
			a.EqStr("entity-2", values[0].ID)
			a.EqStr("entity-2 description", values[0].Desc)

			keys, values, err = testClient.GetAll(ctx, NewEntityQuery().EqID("entity-2").ViaKeys())
			a.Nil(err)
			a.EqInt(1, len(keys))
			a.EqStr("entity-2", values[0].ID)
			a.EqStr("entity-2 description", values[0].Desc)

			var cached = make([]*Entity, 1, 1)
			err = testEnvironment.GetCache().GetMulti(ctx, []string{
				datastore.GetCacheKey(keys[0]),
			}, cached)
			a.Nil(err)
			a.EqStr(values[0].ID, cached[0].ID)
			a.EqStr(values[0].Desc, cached[0].Desc)
		})

		t.Run("Run", func(t *testing.T) {
			a := assert.New(t)
			a.Nil(testEnvironment.Reset())
			a.Nil(testEnvironment.LoadFixture("./fixture/TestEntity_GetAll.json"))
			iter, err := testClient.Run(ctx, NewEntityQuery().EqID("entity-2"))
			a.Nil(err)
			_, value := iter.MustNext()
			a.Nil(err)
			a.EqStr("entity-2", value.ID)
			a.EqStr("entity-2 description", value.Desc)
			_, value = iter.MustNext()
			a.Nil(value)

			iter, err = testClient.Run(ctx, NewEntityQuery().EqID("entity-2").ViaKeys())
			a.Nil(err)
			key, value := iter.MustNext()
			a.Nil(err)
			a.EqStr("entity-2", value.ID)
			a.EqStr("entity-2 description", value.Desc)

			var cached = make([]*Entity, 1, 1)
			err = testEnvironment.GetCache().GetMulti(ctx, []string{
				datastore.GetCacheKey(key),
			}, cached)
			a.Nil(err)
			a.EqStr(value.ID, cached[0].ID)
			a.EqStr(value.Desc, cached[0].Desc)
		})

		t.Run("RunAll", func(t *testing.T) {
			a := assert.New(t)
			a.Nil(testEnvironment.Reset())
			a.Nil(testEnvironment.LoadFixture("./fixture/TestEntity_GetAll.json"))

			_, values, next := testClient.MustRunAll(ctx, NewEntityQuery().EqID("entity-2").ViaKeys())
			a.EqInt(1, len(values))
			a.EqStr("entity-2", values[0].ID)
			a.EqStr("entity-2 description", values[0].Desc)
			a.EqStr("Ci0SJ2oPdGVzdGVudmlyb25tZW50chQLEgZFbnRpdHkiCGVudGl0eS0yDBgAIAA", next)
		})

		t.Run("DeleteMatched", func(t *testing.T) {
			a := assert.New(t)
			a.Nil(testEnvironment.Reset())
			a.Nil(testEnvironment.LoadFixture("./fixture/TestEntity_DeleteMatched.json"))

			a.EqInt(4, testClient.MustCount(ctx, NewEntityQuery()))
			q := NewEntityQuery().LeDigit(2)
			a.EqInt(2, testClient.MustCount(ctx, q))
			deleted := testClient.MustDeleteMatched(ctx, q)
			a.EqInt(2, len(deleted))
			a.EqInt(2, testClient.MustCount(ctx, NewEntityQuery()))
			_, value := testClient.MustGet(ctx, "entity-1")
			a.Nil(value)
			_, value = testClient.MustGet(ctx, "entity-2")
			a.Nil(value)
			_, value = testClient.MustGet(ctx, "entity-3")
			a.NotNil(value)
			_, value = testClient.MustGet(ctx, "entity-4")
			a.NotNil(value)
		})
	})
}
