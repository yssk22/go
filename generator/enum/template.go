package enum

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

{{range .Specs -}}
var (
	_{{.EnumName}}ValueToString = map[{{.EnumName}}]string{
		{{range .Values -}}
		{{.Name}}: "{{.StrValue}}",
		{{end -}}
	}
	_{{.EnumName}}StringToValue = map[string]{{.EnumName}}{
		{{range .Values -}}
		"{{.StrValue}}": {{.Name}},
		{{end -}}
	}
)

func (e {{.EnumName}}) String() string {
	if str, ok := _{{.EnumName}}ValueToString[i]; ok {
		return str
	}
	return fmt.Sprintf("{{.EnumName}}(%d)", e)
}

func Parse{{.EnumName}}(s string) ({{.EnumName}}, error) {
	if val, ok := _{{.EnumName}}StringToValue[s]; ok {
		return val, nil
	}
	return {{.EnumName}}(0), fmt.Errorf("invalid value %q for {{.EnumName}}", s)
}

func Parse{{.EnumName}}Or(s string, or {{.EnumName}}) {{.EnumName}} {
	val, err := Parse{{.EnumName}}(s)
	if err != nil {
		return or
	}
	return val
}

func MustParse{{.EnumName}}(s string) {{.EnumName}} {
	val, err := Parse{{.EnumName}}(s)
	if err != nil {
		panic(err)
	}
	return val
}

func (i {{.EnumName}}) MarshalJSON() ([]byte, error) {
	var s string
	var ok bool
	if s, ok = _{{.EnumName}}ValueToString[i]; !ok {
		s = fmt.Sprintf("{{.EnumName}}(%d)", i)
	}
	return json.Marshal(s)
}

func (i *{{.EnumName}}) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("invalid JSON string")
	}
	newval, err := Parse{{.EnumName}}(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*i = newval
	return nil
}
{{end}}
`
