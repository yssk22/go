package api

const filename = "api.go"

const templateFile = `
package {{.Package}}

import (
    {{range $key, $as := .Dependencies -}}
    {{if $as -}}
    {{$as}} "{{$key}}"
    {{else -}}
    "{{$key}}"
    {{end -}}
    {{end }}
)

func SetupAPI(r *web.Router) {
	{{range .Specs -}}
		r.{{.Method}}("{{.PathPattern}}", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
			obj, err := {{.FuncName}}(req);
			if err != nil {
				return response.NewJSON(err)
			}
			return response.NewJSON(obj)
		}))
	{{end -}}
}
`
