package response

import (
	"net/http/httptest"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
	"golang.org/x/net/context"
)

func TestRedirect(t *testing.T) {
	a := assert.New(t)
	r := NewRedirect("/")
	w := httptest.NewRecorder()
	r.Render(context.Background(), w)
	a.EqInt(int(HTTPStatusSeeOther), w.Code)
}
