package flow

const templateFile = `
{{.Dependency.GenImportForJavaScript}}

{{range .Specs -}}
export type {{.TypeName }} = {{.FlowType.GetExpr}}
{{end -}}
`
