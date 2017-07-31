package response

import (
	"net/http/httptest"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
	"context"
)

func TestJSON(t *testing.T) {
	a := assert.New(t)
	json := NewJSON(map[string]bool{
		"ok": true,
	})
	w := httptest.NewRecorder()
	json.Render(context.Background(), w)

	a.EqStr("{\"ok\":true}\n", w.Body.String())
}

func TestJSON_emptySlice(t *testing.T) {
	a := assert.New(t)
	var list []int
	json := NewJSON(list)
	w := httptest.NewRecorder()
	json.Render(context.Background(), w)

	a.EqStr("[]\n", w.Body.String())
}
