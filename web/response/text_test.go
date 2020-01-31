package response

import (
	"net/http/httptest"
	"testing"

	"context"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestText(t *testing.T) {
	a := assert.New(t)
	text := NewText(context.Background(), "Test Test")
	w := httptest.NewRecorder()
	text.Render(w)

	a.EqStr("Test Test", w.Body.String())
}

func TestTextWithCode(t *testing.T) {
	a := assert.New(t)
	text := NewTextWithStatus(context.Background(), "Test Test", HTTPStatusNotFound)
	w := httptest.NewRecorder()
	text.Render(w)

	a.EqStr("Test Test", w.Body.String())
	a.EqInt(404, w.Code)
}

func TestText_nil(t *testing.T) {
	a := assert.New(t)
	text := NewText(context.Background(), nil)
	w := httptest.NewRecorder()
	text.Render(w)

	a.EqStr("<nil>", w.Body.String())
}
