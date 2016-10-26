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

func TestSession_Encode_Decode(t *testing.T) {
	a := assert.New(t)
	s := NewSession()
	s.Data.Set("FOO", "BAR")
	encoded, err := s.Encode()
	a.Nil(err)

	s1 := NewSession()
	a.Nil(s1.Decode(encoded))
	val, _ := s1.Data.Get("FOO")
	a.EqStr("BAR", val.(string))
}
