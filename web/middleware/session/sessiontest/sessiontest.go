// Package sessiontest provides session test helper
package sessiontest

import (
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/speedland/go/uuid"
	"github.com/speedland/go/web/middleware/session"
	"github.com/speedland/go/x/xnet/xhttp/xhttptest"
	"golang.org/x/net/context"
)

func GetSession(w *httptest.ResponseRecorder, middleware *session.Middleware) (*session.Session, error) {
	c, err := xhttptest.GetCookie(w, middleware.CookieName)
	if err != nil {
		return nil, fmt.Errorf("cookie not found")
	}
	sid, ok := uuid.FromString(strings.Split(c.Value, ".")[0])
	if !ok {
		return nil, fmt.Errorf("cookie does not contain the valid session id")
	}

	return middleware.Store.Get(context.Background(), sid)
}
