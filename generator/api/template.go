package api

const templateFile = `
package {{.Package}}

{{.Dependency.GenImport}}

func SetupAPI(r *web.Router) {
	{{range .Specs -}}
	{{if .StructuredParameter -}}
	var _{{.FuncName}}ParameterParser api.ParameterParser
	json.Unmarshal(
		{{serialize .StructuredParameter.Parser}},
		&_{{.FuncName}}ParameterParser,
	)
	{{end -}}
	r.{{.Method}}("{{.PathPattern}}", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
			{{if .StructuredParameter -}}
			var sp {{.StructuredParameter.Type}}
			if err := _{{.FuncName}}ParameterParser.Parse(req.Request, &sp); err != nil {
				return err.ToResponse()
			}
			{{end -}}
			ctx := req.Context()
			obj, err := {{.FuncName}}(
				ctx,
				{{range .PathParameters -}}
				req.Params.GetStringOr("{{.}}", ""),
				{{end -}}
				{{with .StructuredParameter -}}
				&sp,
				{{end -}}
				);
			if err != nil {
				return api.NewErrorResponse(err)
			}
			return response.NewJSON(obj)
		}))
	{{end -}}
}
`
