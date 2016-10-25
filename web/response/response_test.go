package response

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/speedland/go/x/xcrypto/xhmac"
	"github.com/speedland/go/x/xtesting/assert"
)

func TestSetCookie(t *testing.T) {
	a := assert.New(t)
	hmac := xhmac.NewBase64([]byte("speedland"), nil)
	text := NewText("Test Test")
	c := &http.Cookie{
		Name:  "foo",
		Value: "bar",
	}
	text.SetCookie(c, hmac)

	w := httptest.NewRecorder()
	text.Render(context.Background(), w)

	a.EqStr(
		fmt.Sprintf("foo=%s", hmac.SignString(c.Value)),
		w.Header().Get("set-cookie"),
	)
}
