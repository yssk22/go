package session

import (
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestContextKey(t *testing.T) {
	a := assert.New(t)
	a.EqStr(
		"sessionkey@github.com/speedland/go/web/middleware/session",
		contextKey.String(),
	)
}

func TestSession_Set_Get(t *testing.T) {
	type custom struct {
		Str string
	}
	a := assert.New(t)
	s := NewSession()
	s.Set("FOO", "BAR")
	s.Set("CUSTOM", &custom{
		Str: "custom",
	})
	var str string
	var c custom
	a.Nil(s.Get("FOO", &str))
	a.EqStr("BAR", str)
	a.Nil(s.Get("CUSTOM", &c))
	a.EqStr("custom", c.Str)
}

func TestSession_Encode_Decode(t *testing.T) {
	a := assert.New(t)
	s := NewSession()
	s.Set("FOO", "BAR")
	encoded, err := s.Encode()
	a.Nil(err)

	s1 := NewSession()
	a.Nil(s1.Decode(encoded))
	var str string
	a.Nil(s1.Get("FOO", &str))
	a.EqStr("BAR", str)
}
