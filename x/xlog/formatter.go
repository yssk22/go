package xlog

import (
	"bytes"
	"text/template"
	"time"
)

// Formatter is an inerface to convert *Record to []byte
type Formatter interface {
	Format(*Record) ([]byte, error)
}

// TextFormatter is an implementation of Formatter using text.Template
type TextFormatter struct {
	t       *template.Template
	newline string
}

// NewTextFormatter returns a new *TextFormatter by given template string with default funcMap.
// You can define `t` as a normal "text/template".Template with *Record data. For example
//
//     NewTextFormatter("{{.SourceLine}}")
//
// formats *Record to *Record.SourceLine.
//
// The `Data`` field is formatted by using String() function
//
// Following functions are available in the template and you can add/override custom functions by `NewTextFormatterWithFuncs``
//
//     - formattimestamp
//       function to format .Timestamp field to string. By default, time.RFC3339 format is used.
//
func NewTextFormatter(t string) *TextFormatter {
	return NewTextFormatterWithFuncs(t, nil)
}

// NewTextFormatterWithFuncs is like NewTextFormatter with adding/overriding custom funcions in the template.
func NewTextFormatterWithFuncs(s string, funcMap template.FuncMap) *TextFormatter {
	t := template.New(textFormatterTemplateName)
	if funcMap == nil {
		t = t.Funcs(defaultFuncMap)
	} else {
		// Add defaults to funcMap
		for k, v := range defaultFuncMap {
			if _, ok := funcMap[k]; !ok {
				funcMap[k] = v
			}
		}
		t = t.Funcs(funcMap)
	}
	return &TextFormatter{
		t:       template.Must(t.Parse(s)),
		newline: textFormatterNewline,
	}
}

// Format implements Formatter#Format
func (f *TextFormatter) Format(r *Record) ([]byte, error) {
	var buff bytes.Buffer
	err := f.t.Execute(&buff, r)
	if err != nil {
		return nil, err
	}
	if f.newline != "" {
		_, err = buff.WriteRune('\n')
		if err != nil {
			return nil, err
		}
	}
	return buff.Bytes(), nil
}

const textFormatterTemplateName = "github.com/speedland/go/x/xlog" // template name for TextFormatter
const textFormatterNewline = "\n"

var defaultFuncMap = map[string]interface{}{
	"formattimestamp": func(t time.Time) string {
		return t.Format(time.RFC3339)
	},
}
