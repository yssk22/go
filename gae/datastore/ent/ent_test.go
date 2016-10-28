package ent

import (
	"os"
	"testing"

	"github.com/speedland/go/gae/datastore"
	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/x/xtesting/assert"
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
