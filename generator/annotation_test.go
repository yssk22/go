package generator

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_parseSignatureParams(t *testing.T) {
	a := assert.New(t)
	p := parseSignatureParams("key1=value1 key2=value2")
	a.EqInt(2, len(p))
	a.EqStr("value1", p["key1"])
	a.EqStr("value2", p["key2"])

	p = parseSignatureParams("key1=value1 key2")
	a.EqInt(2, len(p))
	a.EqStr("value1", p["key1"])
	a.EqStr("", p["key2"])

	p = parseSignatureParams("key1=value1key2=fo")
	a.EqInt(1, len(p))
	a.EqStr("value1key2=fo", p["key1"])
}
