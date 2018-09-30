package datastore

import (
	"os"
	"testing"

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
	tt := make([]*Example, 1, 1)
	key := NewKey(gaetest.NewContext(), "Example", "example-1")
	GetMulti(gaetest.NewContext(), []*datastore.Key{key}, tt)
	t.Log(tt[0])
}