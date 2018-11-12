package types

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestTyped(t *testing.T) {
	a := assert.New(t)

	type TestStruct struct {
		Foo string `json:"foo"`
	}
	var v = map[string]interface{}{
		"foo": "bar",
	}
	var vv TestStruct
	Typed(v, &vv)
	a.EqStr("bar", vv.Foo)
}
