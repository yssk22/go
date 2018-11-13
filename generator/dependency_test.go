package generator

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_Dependency(t *testing.T) {
	a := assert.New(t)
	d := NewDependency()
	alias := d.Add("github.com/yssk22/go")
	a.EqStr("go", alias)
	alias = d.Add("github.com/yssk22/go")
	a.EqStr("go", alias)
	alias = d.Add("github.com/foo/go")
	a.EqStr("go1", alias)
	alias = d.Add("github.com/bar/go")
	a.EqStr("go2", alias)

	expect := "import (\n" +
		"\t\"github.com/yssk22/go\"\n" +
		"\tgo1 \"github.com/foo/go\"\n" +
		"\tgo2 \"github.com/bar/go\"\n" +
		")\n"
	a.EqStr(expect, d.GenImport())
}

func Test_Dependency_AddAs(t *testing.T) {
	a := assert.New(t)
	d := NewDependency()
	alias := d.AddAs("github.com/yssk22/go", "foo")
	a.EqStr("foo", alias)
}
