package response

import (
	"net/http/httptest"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestHeader(t *testing.T) {
	a := assert.New(t)
	header := NewHeader()
	header.Set("X-SOME-HEADER", "valuevalue")
	w := httptest.NewRecorder()
	header.Render(w)

	a.EqStr("valuevalue", w.Header().Get("x-some-header"))
}
