package response

import (
	"html/template"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestHTML(t *testing.T) {
	a := assert.New(t)

	var tmpl = template.New("foo")

	template.Must(tmpl.Parse("Sub: {{template \"sub\" .}}"))
	template.Must(tmpl.New("sub").Parse("This is sub {{.foo}}"))

	html := NewHTML(tmpl, map[string]string{
		"foo": "bar",
	})
	w := httptest.NewRecorder()
	html.Render(context.Background(), w)

	a.EqStr("Sub: This is sub bar", w.Body.String())
}
