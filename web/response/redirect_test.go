package response

import (
	"net/http/httptest"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
	"context"
)

func TestRedirect(t *testing.T) {
	a := assert.New(t)
	r := NewRedirect("/")
	w := httptest.NewRecorder()
	r.Render(context.Background(), w)
	a.EqInt(int(HTTPStatusSeeOther), w.Code)
	a.EqStr("/", w.Header().Get("Location"))
}
