package flow

import "github.com/yssk22/go/generator"

type bindings struct {
	Package    string
	Dependency *generator.Dependency
	Specs      []*Spec
}

const templateFile = `
{{.Dependency.GenImportForJavaScript}}

{{range .Specs -}}
export type {{.TypeName }} = {{.FlowType.GetExpr}}
{{end -}}
`
