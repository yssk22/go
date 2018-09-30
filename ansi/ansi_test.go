package ansi

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestColor(t *testing.T) {
	a := assert.New(t)
	s := Blue.Sprintf("  Foo  ")
	a.EqStr("\x1b[34m  Foo  \x1b[0m", s)
}
