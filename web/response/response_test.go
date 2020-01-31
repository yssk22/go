package response

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/yssk22/go/x/xcrypto/xhmac"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestSetCookie(t *testing.T) {
	a := assert.New(t)
	hmac := xhmac.NewBase64([]byte("github.com/yssk22"), nil)
	text := NewText(context.Background(), "Test Test")
	c := &http.Cookie{
		Name:  "foo",
		Value: "bar",
	}
	text.SetCookie(c, hmac)

	w := httptest.NewRecorder()
	text.Render(w)

	a.EqStr(
		fmt.Sprintf("foo=%s", hmac.SignString(c.Value)),
		w.Header().Get("set-cookie"),
	)
}
