package response

import (
	"html/template"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestHTML(t *testing.T) {
	a := assert.New(t)

	var tmpl = template.New("foo")
	var data = make(map[string]interface{})
	tmpl.Parse("Sub: {{sub}}")
	tmpl.New("sub").Parse("This is sub {{.foo}}")
	data["foo"] = "bar"

	html := NewHTML(tmpl, data)
	w := httptest.NewRecorder()
	html.Render(context.Background(), w)

	a.EqStr("Sub: This is sub bar", w.Body.String())
}
