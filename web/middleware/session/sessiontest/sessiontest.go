// Package sessiontest provides session test helper
package sessiontest

import (
	"fmt"
	"strings"

	"context"

	"github.com/yssk22/go/uuid"
	"github.com/yssk22/go/web/middleware/session"
	"github.com/yssk22/go/x/xnet/xhttp/xhttptest"
)

func GetSession(ctx context.Context, w xhttptest.Response, middleware *session.Middleware) (*session.Session, error) {
	c, err := xhttptest.GetCookie(w, middleware.CookieName)
	if err != nil {
		return nil, fmt.Errorf("cookie not found")
	}
	sid, ok := uuid.FromString(strings.Split(c.Value, ".")[0])
	if !ok {
		return nil, fmt.Errorf("cookie does not contain the valid session id")
	}

	return middleware.Store.Get(ctx, sid)
}
