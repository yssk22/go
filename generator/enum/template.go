package enum

import "github.com/yssk22/go/generator"

type bindings struct {
	Package    string
	Dependency *generator.Dependency
	Specs      []*Spec
}

const templateFile = `
package {{.Package}}

{{.Dependency.GenImport}}

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
	if str, ok := _{{.EnumName}}ValueToString[e]; ok {
		return str
	}
	return fmt.Sprintf("{{.EnumName}}(%d)", e)
}

func (e {{.EnumName}}) IsVaild() bool {
	_, ok := _{{.EnumName}}ValueToString[e]
	return ok
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

func (e {{.EnumName}}) MarshalJSON() ([]byte, error) {
	var s string
	var ok bool
	if s, ok = _{{.EnumName}}ValueToString[e]; !ok {
		s = fmt.Sprintf("{{.EnumName}}(%d)", e)
	}
	return json.Marshal(s)
}

func (e *{{.EnumName}}) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("invalid JSON string")
	}
	newval, err := Parse{{.EnumName}}(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*e = newval
	return nil
}

func (e *{{.EnumName}}) Parse(s string) error {
	if val, ok := _{{.EnumName}}StringToValue[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid value %q for {{.EnumName}}", s)
}
{{end}}
`
