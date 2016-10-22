package response

import (
	"net/http/httptest"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
	"golang.org/x/net/context"
)

func TestHeader(t *testing.T) {
	a := assert.New(t)
	header := NewHeader()
	header.Fields.Set("X-SOME-HEADER", "valuevalue")
	w := httptest.NewRecorder()
	header.Render(context.Background(), w)

	a.EqStr("valuevalue", w.Header().Get("x-some-header"))
}

func TestSetHeader(t *testing.T) {
	a := assert.New(t)
	header := NewHeader()

	ctx := context.Background()
	ctx = SetHeader(ctx, "X-SOME-HEADER", "valuevalue")
	w := httptest.NewRecorder()
	header.Render(ctx, w)

	a.EqStr("valuevalue", w.Header().Get("x-some-header"))
}
