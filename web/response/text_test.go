package response

import (
	"net/http/httptest"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
	"context"
)

func TestText(t *testing.T) {
	a := assert.New(t)
	text := NewText("Test Test")
	w := httptest.NewRecorder()
	text.Render(context.Background(), w)

	a.EqStr("Test Test", w.Body.String())
}

func TestTextWithCode(t *testing.T) {
	a := assert.New(t)
	text := NewTextWithStatus("Test Test", HTTPStatusNotFound)
	w := httptest.NewRecorder()
	text.Render(context.Background(), w)

	a.EqStr("Test Test", w.Body.String())
	a.EqInt(404, w.Code)
}

func TestText_nil(t *testing.T) {
	a := assert.New(t)
	text := NewText(nil)
	w := httptest.NewRecorder()
	text.Render(context.Background(), w)

	a.EqStr("<nil>", w.Body.String())
}
