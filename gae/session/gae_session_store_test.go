package session

import (
	"os"
	"testing"

	"github.com/yssk22/go/gae/gaetest"
	"github.com/yssk22/go/web/middleware/session"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestMain(m *testing.M) {
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestGAESessionStore(t *testing.T) {
	type custom struct {
		Str string
	}
	a := assert.New(t)
	store := NewGAESessionStore("")
	sess := session.NewSession()
	sess.Set("foo", "bar")
	sess.Set("custom", &custom{
		Str: "string",
	})
	a.Nil(store.Set(gaetest.NewContext(), sess))

	sess2, err := store.Get(gaetest.NewContext(), sess.ID)
	a.Nil(err)
	a.EqStr(sess.ID.String(), sess2.ID.String())
	var bar string
	var c custom
	a.Nil(sess2.Get("foo", &bar))
	a.EqStr("bar", bar)
	a.Nil(sess2.Get("custom", &c))
	a.EqStr("string", c.Str)
}
