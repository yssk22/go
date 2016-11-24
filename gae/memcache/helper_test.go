package memcache

import (
	"testing"
	"time"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/x/xtesting/assert"
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
