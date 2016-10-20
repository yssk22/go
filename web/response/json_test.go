package response

import (
	"net/http/httptest"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestJSON(t *testing.T) {
	a := assert.New(t)
	json := NewJSON(map[string]bool{
		"ok": true,
	})
	w := httptest.NewRecorder()
	json.Render(w)

	a.EqStr("{\"ok\":true}\n", w.Body.String())
}
