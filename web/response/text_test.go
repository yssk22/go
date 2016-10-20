package response

import (
	"net/http/httptest"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestText(t *testing.T) {
	a := assert.New(t)
	text := NewText("Test Test")
	w := httptest.NewRecorder()
	text.Render(w)

	a.EqStr("Test Test", w.Body.String())
}

func TestTextWithCode(t *testing.T) {
	a := assert.New(t)
	text := NewTextWithCode("Test Test", HTTPStatusNotFound)
	w := httptest.NewRecorder()
	text.Render(w)

	a.EqStr("Test Test", w.Body.String())
	a.EqInt(404, w.Code)
}
