package response

import (
	"net/http/httptest"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
	"golang.org/x/net/context"
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
