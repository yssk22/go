package generator

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_parseAnnotation(t *testing.T) {
	a := assert.New(t)
	ann := parseAnnotation("@foo key1=value1 key2=value2")
	a.EqStr("foo", string(ann.Symbol))
	a.EqInt(2, len(ann.Params))
	a.EqStr("value1", ann.Params["key1"].(string))
	a.EqStr("value2", ann.Params["key2"].(string))

	ann = parseAnnotation("@foo key1=value1 key2")
	a.EqStr("foo", string(ann.Symbol))
	a.EqInt(2, len(ann.Params))
	a.EqStr("value1", ann.Params["key1"].(string))
	a.Nil(ann.Params["key2"])

	ann = parseAnnotation("@bar key1=value1key2=fo")
	a.EqStr("bar", string(ann.Symbol))
	a.EqInt(1, len(ann.Params))
	a.EqStr("value1key2=fo", ann.Params["key1"].(string))
}
