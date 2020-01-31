package response

import (
	"net/http/httptest"
	"testing"

	"context"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestJSON(t *testing.T) {
	a := assert.New(t)
	json := NewJSON(context.Background(), map[string]bool{
		"ok": true,
	})
	w := httptest.NewRecorder()
	json.Render(w)

	a.EqStr("{\"ok\":true}\n", w.Body.String())
}

func TestJSON_emptySlice(t *testing.T) {
	a := assert.New(t)
	var list []int
	json := NewJSON(context.Background(), list)
	w := httptest.NewRecorder()
	json.Render(w)

	a.EqStr("[]\n", w.Body.String())
}
