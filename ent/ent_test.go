package ent

import (
	"os"
	"testing"

	"github.com/yssk22/go/gae/datastore"
	"github.com/yssk22/go/gae/gaetest"
	"github.com/yssk22/go/x/xtesting/assert"
	"google.golang.org/appengine"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestGetMemcacheKey(t *testing.T) {
	a := assert.New(t)
	k := datastore.NewKey(gaetest.NewContext(), "MyKind", "MyKey")
	a.EqStr("datastore.MyKind.MyKey", GetMemcacheKey(k))

	ns, _ := appengine.Namespace(gaetest.NewContext(), "ns")
	k = datastore.NewKey(
		ns,
		"MyKind", "MyKey",
	)
	a.EqStr("datastore.MyKind.MyKey", GetMemcacheKey(k))
}
