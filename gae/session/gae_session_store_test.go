package session

import (
	"os"
	"testing"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/web/middleware/session"
	"github.com/speedland/go/x/xtesting/assert"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestGAESessionStore(t *testing.T) {
	a := assert.New(t)
	store := NewGAESessionStore("")
	sess := session.NewSession()
	sess.Data.Set("foo", "bar")
	a.Nil(store.Set(gaetest.NewContext(), sess))

	sess2, err := store.Get(gaetest.NewContext(), sess.ID)
	a.Nil(err)
	a.EqStr(sess.ID.String(), sess2.ID.String())
	bar, err := sess2.Data.Get("foo")
	a.Nil(err)
	a.EqStr("bar", bar.(string))
}
