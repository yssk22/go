package datastore

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	c := NewClientFromClient(ctx, testEnv.client, Cache(testEnv.memcache))

	t.Run("GetMulti", func(t *testing.T) {
		a := assert.New(t)
		a.Nil(testEnv.Reset())
		a.Nil(testEnv.LoadFixture("./fixtures/TestGetMulti.json"))

		tt := make([]*Example, 2, 2)
		keys := []*datastore.Key{
			NewKey(ctx, "Example", "example-1"),
			NewKey(ctx, "Example", "example-3"),
		}
		a.Nil(c.GetMulti(ctx, keys, tt))
		a.NotNil(tt[0])
		a.Nil(tt[1])

		e1 := make([]*Example, 1)
		err := testEnv.memcache.GetMulti(ctx, []string{
			GetCacheKey(keys[0]),
		}, e1)
		a.Nil(err)
		a.EqStr("example-1", e1[0].ID)

		keys = []*datastore.Key{
			NewKey(ctx, "Example", "example-1"),
			NewKey(ctx, "Example", "example-3"),
			NewKey(ctx, "Example", "example-2"),
		}
		tt = make([]*Example, 3, 3)
		a.Nil(c.GetMulti(ctx, keys, tt))
		a.NotNil(tt[0])
		a.Nil(tt[1])
		a.NotNil(tt[2])
	})

	t.Run("PutMulti", func(t *testing.T) {
		a := assert.New(t)
		a.Nil(testEnv.Reset())
		ctx := context.Background()
		c := NewClientFromClient(ctx, testEnv.client, Cache(testEnv.memcache))
		keys := []*datastore.Key{
			NewKey(ctx, "Example", "example-a"),
		}
		stored := make([]*Example, 1, 1)
		a.Nil(c.GetMulti(ctx, keys, stored))

		_, err := c.PutMulti(ctx, keys, []*Example{
			{
				ID: "example-a",
			},
		})
		a.Nil(err)

		caches := make([]*Example, 1, 1)
		a.NotNil(testEnv.memcache.GetMulti(ctx, []string{GetCacheKey(keys[0])}, caches))
		stored = make([]*Example, 1, 1)
		a.Nil(testEnv.client.GetMulti(ctx, keys, stored))
		a.NotNil(stored[0])
		a.EqStr("example-a", stored[0].ID)
	})

	t.Run("DeleteMulti", func(t *testing.T) {
		a := assert.New(t)
		a.Nil(testEnv.Reset())
		a.Nil(testEnv.LoadFixture("./fixtures/TestGetMulti.json"))

		ctx := context.Background()
		c := NewClientFromClient(ctx, testEnv.client, Cache(testEnv.memcache))

		keys := []*datastore.Key{
			NewKey(ctx, "Example", "example-1"),
			NewKey(ctx, "Example", "example-3"),
		}

		a.Nil(c.DeleteMulti(ctx, keys))

		tt := make([]*Example, 2, 2)
		a.Nil(c.GetMulti(ctx, keys, tt))
		a.Nil(tt[0])
		a.Nil(tt[1])
	})

	t.Run("Query", func(t *testing.T) {
		a := assert.New(t)
		a.Nil(testEnv.Reset())
		a.Nil(testEnv.LoadFixture("./fixtures/TestQuery.json"))

		t.Run("KeysOnly", func(t *testing.T) {
			a := assert.New(t)
			q := NewQuery("Example").Eq("ID", "example-1").KeysOnly()
			keys, err := c.GetAll(ctx, q, nil)
			a.Nil(err)
			a.EqInt(1, len(keys))
		})

		t.Run("Filter", func(t *testing.T) {
			a := assert.New(t)
			var result []Example
			q := NewQuery("Example").Eq("ID", "example-1")
			_, err := c.GetAll(ctx, q, &result)
			a.Nil(err)
			a.EqInt(1, len(result))
		})

		t.Run("Order", func(t *testing.T) {
			a := assert.New(t)
			var result []Example
			q := NewQuery("Example").Desc("ID")
			_, err := c.GetAll(ctx, q, &result)
			a.Nil(err)
			a.EqInt(5, len(result))
			for i := range result {
				a.EqStr(fmt.Sprintf("example-%d", 5-i), result[i].ID)
			}
		})

		t.Run("ViaKeys", func(t *testing.T) {

		})
		t.Run("Limit", func(t *testing.T) {
			a := assert.New(t)
			var result []Example
			q := NewQuery("Example").Desc("ID").Limit(1)
			_, err := c.GetAll(ctx, q, &result)
			a.Nil(err)
			a.EqInt(1, len(result))
			a.EqStr("example-5", result[0].ID)
		})
	})
}
