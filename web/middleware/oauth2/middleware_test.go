package oauth2

import (
	"fmt"
	"testing"

	"github.com/speedland/go/uuid"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/httptest"
	"github.com/speedland/go/web/middleware/session"
	"github.com/speedland/go/web/middleware/session/sessiontest"
	"github.com/speedland/go/web/response"

	"golang.org/x/oauth2"
)

var sessionMiddleware = session.NewMiddleware()

func TestMiddleware_Redirect(t *testing.T) {
	middleware := &Middleware{}
	middleware.AuthPath = "/oauth2/login"
	middleware.Store = &SessionStateStore{}
	middleware.Config = &TestConfig{}

	a := httptest.NewAssert(t)
	router := prepareRouter(middleware)

	res := router.TestGet("/oauth2/login")
	a.Status(response.HTTPStatusFound, res)
	session, err := sessiontest.GetSession(res, sessionMiddleware)
	a.Nil(err)

	v, err := session.Get(oauth2SessionStateKey)
	a.Nil(err)
	uuid, ok := uuid.FromString(v.(string))
	a.OK(ok)
	a.Header(
		fmt.Sprintf("http://oauth2.example.com/?state=%s", uuid.String()),
		res, "Location",
	)
}

func prepareRouter(middleware *Middleware) *httptest.Router {
	router := httptest.NewRouter(nil)
	router.Use(sessionMiddleware)
	router.Use(middleware)
	router.Get("/", web.HandlerFunc(func(req *web.Request, _ web.NextHandler) *response.Response {
		return response.NewText("OK")
	}))
	return router
}

type TestConfig struct{}

func (c *TestConfig) AuthCodeURL(state string, _ ...oauth2.AuthCodeOption) string {
	return fmt.Sprintf("http://oauth2.example.com/?state=%s", state)
}
