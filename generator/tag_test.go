package generator

import (
	"testing"

	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_Tag(t *testing.T) {
	a := assert.New(t)
	tags := ParseTag(`json:"foo" value:"1"`)
	a.EqStr("foo", keyvalue.GetStringOr(tags, "json", ""))
	a.EqInt(1, keyvalue.GetIntOr(tags, "value", 0))
}
