package xlog

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/speedland/go/ansi"
)

// Formatter is an inerface to convert *Record to []byte
type Formatter interface {
	Format(*Record) ([]byte, error)
}

// TextFormatter is an implementation of Formatter using text.Template
type TextFormatter struct {
	d       *textFormatterTemplate // default
	funcMap template.FuncMap
	trace   *textFormatterTemplate
	debug   *textFormatterTemplate
	info    *textFormatterTemplate
	warn    *textFormatterTemplate
	error   *textFormatterTemplate
	fatal   *textFormatterTemplate
	newline []byte
}

// Format immplements Formatter#Format(*Record)
func (f *TextFormatter) Format(r *Record) ([]byte, error) {
	var t *textFormatterTemplate
	switch r.Level {
	case LevelTrace:
		t = f.trace
		break
	case LevelDebug:
		t = f.debug
		break
	case LevelInfo:
		t = f.info
		break
	case LevelWarn:
		t = f.warn
		break
	case LevelError:
		t = f.error
		break
	case LevelFatal:
		t = f.fatal
		break
	}
	if t == nil {
		t = f.d
	}
	buff, err := t.Format(r)
	if err != nil {
		return nil, err
	}
	if f.newline != nil {
		return append(buff, f.newline...), nil
	}
	return buff, nil
}

// SetCode sets the default ANSI code for log texts
func (f *TextFormatter) SetCode(code ansi.Code) *TextFormatter {
	f.d.code = code
	return f
}

// SetTrace set the template for LevelTrace logs.
func (f *TextFormatter) SetTrace(s string, code ansi.Code) *TextFormatter {
	f.trace = newTextFormatterTemplate(s, code, f.funcMap)
	f.trace.code = code
	return f
}

// SetDebug set the template for LevelDebug logs.
func (f *TextFormatter) SetDebug(s string, code ansi.Code) *TextFormatter {
	f.debug = newTextFormatterTemplate(s, code, f.funcMap)
	return f
}

// SetInfo set the template for LevelInfo logs.
func (f *TextFormatter) SetInfo(s string, code ansi.Code) *TextFormatter {
	f.info = newTextFormatterTemplate(s, code, f.funcMap)
	return f
}

// SetWarn set the template for LevelWarn logs.
func (f *TextFormatter) SetWarn(s string, code ansi.Code) *TextFormatter {
	f.warn = newTextFormatterTemplate(s, code, f.funcMap)
	return f
}

// SetError set the template for LevelError logs.
func (f *TextFormatter) SetError(s string, code ansi.Code) *TextFormatter {
	f.error = newTextFormatterTemplate(s, code, f.funcMap)
	return f
}

// SetFatal set the template for LevelFatal logs.
func (f *TextFormatter) SetFatal(s string, code ansi.Code) *TextFormatter {
	f.fatal = newTextFormatterTemplate(s, code, f.funcMap)
	return f
}

type textFormatterTemplate struct {
	t    *template.Template
	code ansi.Code
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
//     - formatstack
//       function to format .Stack field to string. By default, All of stack frames are rendered.
//
func NewTextFormatter(t string) *TextFormatter {
	return NewTextFormatterWithFuncs(t, nil)
}

// NewTextFormatterWithFuncs is like NewTextFormatter with adding/overriding custom funcions in the template.
func NewTextFormatterWithFuncs(s string, funcMap template.FuncMap) *TextFormatter {
	return &TextFormatter{
		d:       newTextFormatterTemplate(s, ansi.Reset, funcMap),
		funcMap: funcMap,
		newline: []byte(textFormatterNewline),
	}
}

func newTextFormatterTemplate(s string, code ansi.Code, funcMap template.FuncMap) *textFormatterTemplate {
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
	t = t.Option("missingkey=zero")
	return &textFormatterTemplate{
		t:    template.Must(t.Parse(s)),
		code: code,
	}

}

// Format implements Formatter#Format
func (f *textFormatterTemplate) Format(r *Record) ([]byte, error) {
	var buff bytes.Buffer
	err := f.t.Execute(&buff, r)
	if err != nil {
		panic(err)
	}
	if f.code == ansi.Reset {
		return buff.Bytes(), nil
	}
	return []byte(f.code.Sprintf(buff.String())), nil
}

const textFormatterTemplateName = "github.com/speedland/go/x/xlog" // template name for TextFormatter
const textFormatterNewline = "\n"

var defaultFuncMap = map[string]interface{}{
	"formattimestamp": func(r *Record) string {
		return r.Timestamp.Format(time.RFC3339)
	},
	"formatstack": func(r *Record) string {
		if r.Stack == nil {
			return "<No stack available>"
		}
		var i = 0
		var buff bytes.Buffer
		for i = range r.Stack {
			if r.Stack[i].PackageName == "" {
				return buff.String()
			}
			buff.WriteString(fmt.Sprintf("\n\t%s", r.Stack[i].String()))
		}
		if i == 50 {
			buff.WriteString("\n\t...and more")
			return buff.String()
		}
		return buff.String()
	},
}
