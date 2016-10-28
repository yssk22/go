package example

import (
	"os"
	"testing"
	"time"

	"github.com/speedland/go/x/xlog"
	"github.com/speedland/go/x/xtime"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/x/xtesting/assert"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
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

func TestExampleKind_GetMulti(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixture/TestExample_GetMulti.json", nil))

	k := &ExampleKind{}
	keys, values, err := k.GetMulti(gaetest.NewContext(), "example-1", "example-2")
	a.Nil(err)
	a.EqInt(2, len(keys))
	a.EqInt(2, len(values))
	a.NotNil(values[0])
	a.NotNil(values[1])
}

func TestExampleKind_GetMulti_notFound(t *testing.T) {
	xlog.SetKeyFilter(ExampleKindLoggerKey, xlog.LevelDebug)
	a := assert.New(t)
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
	xlog.SetKeyFilter(ExampleKindLoggerKey, xlog.LevelDebug)
	a := assert.New(t)
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
