package cache

import (
	"context"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestMemoryCache(t *testing.T) {
	ctx := context.Background()

	t.Run("same type", func(t *testing.T) {
		a := assert.New(t)
		mc := &MemoryCache{}
		a.Nil(mc.SetMulti(ctx, []string{"1"}, []*Example{
			{ID: "1"},
		}))
		cached := make([]*Example, 1, 1)
		a.Nil(mc.GetMulti(ctx, []string{"1"}, cached))
		a.EqInt(1, len(cached))
		a.EqStr("1", cached[0].ID)
	})

	t.Run("ptr to non ptr", func(t *testing.T) {
		a := assert.New(t)
		mc := &MemoryCache{}
		a.Nil(mc.SetMulti(ctx, []string{"1"}, []*Example{
			{ID: "1"},
		}))
		cached := make([]Example, 1, 1)
		a.Nil(mc.GetMulti(ctx, []string{"1"}, cached))
		a.EqInt(1, len(cached))
		a.EqStr("1", cached[0].ID)
	})

	t.Run("non ptr to ptr", func(t *testing.T) {
		a := assert.New(t)
		mc := &MemoryCache{}
		a.Nil(mc.SetMulti(ctx, []string{"1"}, []Example{
			{ID: "1"},
		}))
		cached := make([]*Example, 1, 1)
		a.Nil(mc.GetMulti(ctx, []string{"1"}, cached))
		a.EqInt(1, len(cached))
		a.EqStr("1", cached[0].ID)
	})

	t.Run("invalid type", func(t *testing.T) {
		a := assert.New(t)
		mc := &MemoryCache{}
		a.Nil(mc.SetMulti(ctx, []string{"1"}, []Example{
			{ID: "1"},
		}))
		cached := make([]**Example, 1, 1)
		a.NotNil(mc.GetMulti(ctx, []string{"1"}, cached))
	})

	t.Run("invalid key", func(t *testing.T) {
		a := assert.New(t)
		mc := &MemoryCache{}
		a.Nil(mc.SetMulti(ctx, []string{"1"}, []Example{
			{ID: "1"},
		}))
		cached := make([]Example, 1, 1)
		a.NotNil(mc.GetMulti(ctx, []string{"2"}, cached))
	})

}
