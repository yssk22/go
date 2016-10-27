package session

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/speedland/go/uuid"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/httptest"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xnet/xhttp/xhttptest"
	"github.com/speedland/go/x/xtime"
	"golang.org/x/net/context"
)

func TestMiddleware_NewSession(t *testing.T) {
	middleware := NewMiddleware()
	sessionDataKey := "FOO"
	sessionDataValue := "BAR"

	a := httptest.NewAssert(t)
	router := prepareRouter(sessionDataKey, sessionDataValue, middleware)
	res := router.TestPost("/session", nil)
	a.NotNil(res)
	a.Status(response.HTTPStatusOK, res)
	a.Body("FOO", res)

	// Check the resposne cookie contains `CookieName`
	c, _ := xhttptest.GetCookie(res, middleware.CookieName)
	a.NotNil(c)
	sid, ok := uuid.FromString(strings.Split(c.Value, ".")[0])
	a.OK(ok)

	// Check the store has session.
	session, err := middleware.Store.Get(context.Background(), sid)
	a.Nil(err)
	value, err := session.Get("FOO")
	a.Nil(err)
	a.EqStr(sessionDataValue, value.(string))

	// Make another Request with cookie
	req, _ := http.NewRequest("GET", "/session", nil)
	req.AddCookie(c)
	res = router.TestRequest(req)
	a.Status(response.HTTPStatusOK, res)
	a.Body(sessionDataValue, res)
}

func TestMiddleware_NoSessionCreation(t *testing.T) {
	middleware := NewMiddleware()
	sessionDataKey := "FOO"
	sessionDataValue := "BAR"

	a := httptest.NewAssert(t)
	router := prepareRouter(sessionDataKey, sessionDataValue, middleware)
	res := router.TestGet("/session")
	a.Status(response.HTTPStatusOK, res)
	a.Body("<nil>", res)

	// No SessionData is stored, so SessionStore should not
	store := middleware.Store.(*MemorySessionStore)
	a.EqInt(0, len(store.store))
}

func TestMiddleware_SessionExpiration(t *testing.T) {
	var c *http.Cookie
	middleware := NewMiddleware()
	sessionDataKey := "FOO"
	sessionDataValue := "BAR"
	a := httptest.NewAssert(t)
	router := prepareRouter(sessionDataKey, sessionDataValue, middleware)

	xtime.RunAt(
		time.Date(2015, 1, 1, 0, 0, 0, 0, xtime.JST),
		func() {
			res := router.TestPost("/session", nil)
			a.NotNil(res)
			a.Status(response.HTTPStatusOK, res)
			a.Body("FOO", res)
			c, _ = xhttptest.GetCookie(res, middleware.CookieName)
			a.NotNil(c)
		},
	)

	xtime.RunAt(
		time.Date(2016, 1, 1, 0, 0, 0, 0, xtime.JST),
		func() {
			req := httptest.NewRequest("GET", "/session", nil)
			req.AddCookie(c)
			res := router.TestRequest(req)
			a.NotNil(res)
			a.Status(response.HTTPStatusOK, res)
			a.Body("<nil>", res)
		},
	)

	// Ensure it is deleted
	store := middleware.Store.(*MemorySessionStore)
	a.EqInt(0, len(store.store))
}

func prepareRouter(sessionDataKey, sessionDataValue interface{}, middleware *Middleware) *httptest.Router {
	router := httptest.NewRouter(nil)
	router.Use(middleware)
	router.Post("/session", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		session := FromContext(req.Context())
		session.Set(sessionDataKey, sessionDataValue)
		return response.NewText(sessionDataKey)
	}))
	router.Get("/session", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		session := FromContext(req.Context())
		v, _ := session.Get(sessionDataKey)
		return response.NewText(v)
	}))
	router.Get("/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return response.NewText("ok")
	}))
	return router
}
