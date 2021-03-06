package session

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"context"
	"github.com/yssk22/go/uuid"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/httptest"
	"github.com/yssk22/go/web/response"
	"github.com/yssk22/go/x/xnet/xhttp/xhttptest"
	"github.com/yssk22/go/x/xtime"
)

func TestMiddleware_NewSession(t *testing.T) {
	middleware := NewMiddleware(NewMemorySessionStore())
	sessionDataKey := "FOO"
	sessionDataValue := "BAR"

	a := httptest.NewAssert(t)
	recorder := httptest.NewRecorder(prepareRouter(sessionDataKey, sessionDataValue, middleware))
	res := recorder.TestPost("/session", nil)
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
	var value string
	a.Nil(err)
	a.Nil(session.Get("FOO", &value))
	a.EqStr(sessionDataValue, value)

	// Make another Request on the same recorder
	res = recorder.TestGet("/session")
	a.Status(response.HTTPStatusOK, res)
	a.Body(sessionDataValue, res)
}

func TestMiddleware_NoSessionCreation(t *testing.T) {
	middleware := NewMiddleware(NewMemorySessionStore())
	sessionDataKey := "FOO"
	sessionDataValue := "BAR"

	a := httptest.NewAssert(t)
	recorder := httptest.NewRecorder(prepareRouter(sessionDataKey, sessionDataValue, middleware))
	res := recorder.TestGet("/session")
	a.Status(response.HTTPStatusOK, res)
	a.Body("<nil>", res)

	// No SessionData is stored, so SessionStore should not
	store := middleware.Store.(*MemorySessionStore)
	a.EqInt(0, len(store.store))
}

func TestMiddleware_SessionExpiration(t *testing.T) {
	var c *http.Cookie
	middleware := NewMiddleware(NewMemorySessionStore())
	sessionDataKey := "FOO"
	sessionDataValue := "BAR"
	a := httptest.NewAssert(t)
	recorder := httptest.NewRecorder(prepareRouter(sessionDataKey, sessionDataValue, middleware))

	xtime.RunAt(
		time.Date(2015, 1, 1, 0, 0, 0, 0, xtime.JST),
		func() {
			res := recorder.TestPost("/session", nil)
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
			res := recorder.TestGet("/session")
			a.NotNil(res)
			a.Status(response.HTTPStatusOK, res)
			a.Body("<nil>", res)
		},
	)

	// Ensure it is deleted
	a.EqInt(0, len(middleware.Store.(*MemorySessionStore).store))
}

func prepareRouter(sessionDataKey, sessionDataValue interface{}, middleware *Middleware) web.Router {
	router := web.NewRouter(nil)
	router.Use(middleware)
	router.Post("/session", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		session := FromContext(req.Context())
		session.Set(sessionDataKey, sessionDataValue)
		return response.NewText(req.Context(), sessionDataKey)
	}))
	router.Get("/session", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		session := FromContext(req.Context())
		if err := session.Get(sessionDataKey, &sessionDataValue); err != nil {
			return response.NewText(req.Context(), nil)
		}
		return response.NewText(req.Context(), sessionDataValue)
	}))
	router.Get("/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return response.NewText(req.Context(), "ok")
	}))
	return router
}
