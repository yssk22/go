package datastore

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/yssk22/go/x/xstrings"
)

var whereFuncs = []string{
	"Eq", "Lt", "Le", "Gt", "Ge", "Ne",
}

var orderFuncs = []string{
	"Asc", "Desc",
}

const whereFuncTemplate = `
func (d *%sQuery) %s%s(v %s) *%sQuery {
	d.query = d.query.%s("%s", v)
	return d
}
`

const orderFuncTemplate = `
func (d *%sQuery) %s%s() *%sQuery {
	d.query = d.query.%s("%s")
	return d
}
`

var templateHelper = template.FuncMap(map[string]interface{}{
	"snakecase": func(s string) string {
		return xstrings.ToSnakeCase(s)
	},
	"mkPrivate": func(s string) string {
		// FooBar => fooBar
		return fmt.Sprintf("%s%s", strings.ToUpper(string(s[0])), string(s[0:]))
	},
	"queryFuncs": func(spec *Spec) string {
		// generate EqXXX() like query funcs.
		var funcs []string
		for _, funcName := range whereFuncs {
			for _, querySpec := range spec.QuerySpecs {
				funcs = append(funcs,
					fmt.Sprintf(whereFuncTemplate,
						spec.StructName,
						funcName,
						querySpec.Name,
						querySpec.Type,
						spec.StructName,
						funcName,
						querySpec.PropertyName,
					))
			}
		}
		for _, funcName := range orderFuncs {
			for _, querySpec := range spec.QuerySpecs {
				funcs = append(funcs,
					fmt.Sprintf(orderFuncTemplate,
						spec.StructName,
						funcName,
						querySpec.Name,
						spec.StructName,
						funcName,
						querySpec.PropertyName,
					))
			}
		}
		return strings.Join(funcs, "\n")
	},
})
