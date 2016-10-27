package example

import (
	"os"
	"testing"

	"github.com/speedland/go/x/xlog"

	"github.com/speedland/go/web/gae/gaetest"
	"github.com/speedland/go/x/xtesting/assert"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestExmample_GetMulti(t *testing.T) {
	xlog.SetKeyFilter(ExampleKindLoggerKey, xlog.LevelDebug)
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

func TestExmample_GetMulti_notFound(t *testing.T) {
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
