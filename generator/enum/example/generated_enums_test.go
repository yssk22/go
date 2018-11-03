package example

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestEnum(t *testing.T) {
	a := assert.New(t)
	a.EqStr("a", MyEnumA.String())
	a.OK(MyEnum(1).IsVaild())
	a.OK(!MyEnum(5).IsVaild())
	a.OK(MyEnum(1) == MustParseMyEnum("b"))

	var v MyEnum
	a.EqInt(0, int(v))
	a.Nil(v.Parse("b"))
	a.EqInt(1, int(v))

	a.NotNil(v.Parse("invalid"))
	a.EqInt(1, int(v))
}
