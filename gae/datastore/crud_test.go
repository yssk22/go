package datastore

import (
	"os"
	"testing"

	"github.com/yssk22/go/gae/memcache"
	"github.com/yssk22/go/gae/gaetest"
	"github.com/yssk22/go/x/xtesting/assert"
	"google.golang.org/appengine/datastore"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

type Example struct {
	ID string
}

func TestGetMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestGetMulti.json", nil))
	ctx := gaetest.NewContext()
	tt := make([]*Example, 2, 2)
	keys := []*datastore.Key{
		NewKey(ctx, "Example", "example-1"),
		NewKey(ctx, "Example", "example-3"),
	}

	a.Nil(GetMulti(ctx, keys, tt))
	a.NotNil(tt[0])
	a.Nil(tt[1])

	var e1 Example
	err := memcache.Get(gaetest.NewContext(), GetMemcacheKey(keys[0]), &e1)
	a.Nil(err)
	a.EqStr("example-1", e1.ID)

	keys = []*datastore.Key{
		NewKey(ctx, "Example", "example-1"),
		NewKey(ctx, "Example", "example-3"),
		NewKey(ctx, "Example", "example-2"),
	}
	tt = make([]*Example, 3, 3)
	a.Nil(GetMulti(ctx, keys, tt))
	a.NotNil(tt[0])
	a.Nil(tt[1])
	a.NotNil(tt[2])
}

func TestPutMulti(t *testing.T) {
	a := assert.New(t)
	gaetest.CleanupDatastore(gaetest.NewContext())	
	ctx := gaetest.NewContext()
	keys := []*datastore.Key{
		NewKey(ctx, "Example", "example-a"),
	}
	_, err := PutMulti(ctx, keys, []*Example{
		{
			ID: "example-a",
		},
	})
	a.Nil(err)
	a.OK(!memcache.Exists(ctx, GetMemcacheKey(keys[0])))
	tt := make([]*Example, 1, 1)
	a.Nil(GetMulti(ctx, keys, tt))
	a.NotNil(tt[0])
	a.EqStr("example-a", tt[0].ID)
}

func TestDeleteMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestGetMulti.json", nil))
	ctx := gaetest.NewContext()
	keys := []*datastore.Key{
		NewKey(ctx, "Example", "example-1"),
		NewKey(ctx, "Example", "example-3"),
	}

	a.Nil(DeleteMulti(ctx, keys))

	tt := make([]*Example, 2, 2)
	a.Nil(GetMulti(ctx, keys, tt))
	a.Nil(tt[0])
	a.Nil(tt[1])
}