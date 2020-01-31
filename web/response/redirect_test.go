package response

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestRedirect(t *testing.T) {
	a := assert.New(t)
	r := NewRedirect(context.Background(), "/")
	w := httptest.NewRecorder()
	r.Render(w)
	a.EqInt(int(HTTPStatusSeeOther), w.Code)
	a.EqStr("/", w.Header().Get("Location"))
}
