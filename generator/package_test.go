package generator

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_resolveGoImportPath(t *testing.T) {
	a := assert.New(t)
	resolved, err := resolveGoImportPath("./testdata/testmod")
	a.Nil(err)
	a.EqStr("github.com/yssk22/go/testmod", resolved)

	resolved, err = resolveGoImportPath("./testdata/testmod/mymodule")
	a.Nil(err)
	a.EqStr("github.com/yssk22/go/testmod/mymodule", resolved)
}
