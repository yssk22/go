package generator

import (
	"bytes"
	"go/format"
	"io"
	"os/exec"

	"github.com/yssk22/go/x/xerrors"
)

type Formatter interface {
	Format(string) (string, error)
}

// FormatterFunc is to define Formatter interface from a func
type FormatterFunc func(string) (string, error)

// Format implements Formatter#format
func (f FormatterFunc) Format(s string) (string, error) {
	return f(s)
}

// GoFormatter is a formatter for go.
var GoFormatter = FormatterFunc(func(src string) (string, error) {
	formatted, err := format.Source([]byte(src))
	if err != nil {
		return "", &InvalidSourceError{
			Source: src,
			err:    err,
		}
	}
	return string(formatted), nil
})

// JavaScriptFormatter is a formatter for js
var JavaScriptFormatter = FormatterFunc(func(src string) (string, error) {
	var out bytes.Buffer
	c := exec.Command("prettier")
	r, w := io.Pipe()
	c.Stdin = r
	c.Stdout = &out
	if err := c.Start(); err != nil {
		w.Close()
		return "", xerrors.Wrap(err, "cannot launch prettier")
	}
	if _, err := w.Write([]byte(src)); err != nil {
		w.Close()
		return "", err
	}
	w.Close()
	c.Wait()
	return out.String(), nil
})
