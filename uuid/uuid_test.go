package uuid

import (
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestFromString(t *testing.T) {
	a := assert.New(t)
	s := New()
	s1, ok := FromString(s.String())
	a.OK(ok)
	a.OK(s == s1)
}
