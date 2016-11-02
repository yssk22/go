package oauth2

import (
	"fmt"
	"net/url"
	"testing"

	"golang.org/x/net/context"

	"github.com/speedland/go/uuid"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/httptest"
	"github.com/speedland/go/web/middleware/session"
	"github.com/speedland/go/web/middleware/session/sessiontest"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xnet/xhttp/xhttptest"

	"golang.org/x/oauth2"
)

var sessionMiddleware = session.NewMiddleware()

func TestMiddleware(t *testing.T) {
	middleware := &Middleware{}
	middleware.AuthPath = "/oauth2/login"
	middleware.CallbackPath = "/oauth2/callback"
	middleware.Store = &SessionStateStore{}
	middleware.Config = &TestConfig{}

	a := httptest.NewAssert(t)
	recorder := httptest.NewRecorder(prepareRouter(middleware))

	// 1. Redirect (to prepare auth state key)
	res := recorder.TestGet("/oauth2/login")
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

	// 2. callback
	cookie, _ := xhttptest.GetCookie(res, sessionMiddleware.CookieName)
	req := recorder.NewRequest("POST", "/oauth2/callback", url.Values{
		"code":  []string{"validcode"},
		"state": []string{uuid.String()},
	})
	req.AddCookie(cookie)
	res = recorder.TestRequest(req)
	a.Body("AccessToken: test-access-token, RefreshToken: test-refresh-token", res)
}

func prepareRouter(middleware *Middleware) *web.Router {
	router := web.NewRouter(nil)
	router.Use(sessionMiddleware)
	router.Use(middleware)
	router.Get("/", web.HandlerFunc(func(req *web.Request, _ web.NextHandler) *response.Response {
		return response.NewText("OK")
	}))
	router.Post("/oauth2/callback", web.HandlerFunc(func(req *web.Request, _ web.NextHandler) *response.Response {
		token := FromContext(req.Context())
		if token == nil {
			return response.NewError(fmt.Errorf("token not found"))
		}
		return response.NewText(
			fmt.Sprintf(
				"AccessToken: %s, RefreshToken: %s",
				token.AccessToken, token.RefreshToken,
			),
		)
	}))
	return router
}

type TestConfig struct {
	Err error
}

func (c *TestConfig) AuthCodeURL(state string, _ ...oauth2.AuthCodeOption) string {
	return fmt.Sprintf("http://oauth2.example.com/?state=%s", state)
}

func (c *TestConfig) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	if c.Err != nil {
		return nil, c.Err
	}
	return &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
	}, nil
}
