package counter

import (
	"os"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"

	"github.com/speedland/go/gae/gaetest"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestCounter(t *testing.T) {
	a := assert.New(t)
	const ckey1 = "testcounter"
	const ckey2 = "testcounter2"
	a.Nil(Increment(gaetest.NewContext(), ckey1))
	a.Nil(Increment(gaetest.NewContext(), ckey1))
	a.Nil(Increment(gaetest.NewContext(), ckey1))
	a.Nil(Increment(gaetest.NewContext(), ckey1))
	a.EqInt(4, MustCount(gaetest.NewContext(), ckey1))

	a.Nil(Increment(gaetest.NewContext(), ckey2))
	a.Nil(Increment(gaetest.NewContext(), ckey2))
	a.EqInt(2, MustCount(gaetest.NewContext(), ckey2))
}
