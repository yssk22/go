package oauth2

import (
	"github.com/speedland/go/uuid"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
)

const OAuth2LoggerKey = "web.middleware.oauth2"

type Middleware struct {
	AuthPath     string // path that redirects to oauth2 provider
	CallbackPath string
	Store        StateStore

	Config Config
}

func NewMiddleware() *Middleware {
	m := &Middleware{}
	m.Store = &SessionStateStore{}
	return m
}

func (m *Middleware) Process(req *web.Request, next web.NextHandler) *response.Response {
	if m.AuthPath == req.URL.EscapedPath() {
		return m.handleAuthPath(req)
	}
	return next(req)
}

func (m *Middleware) handleAuthPath(req *web.Request) *response.Response {
	state := uuid.New().String()
	if err := m.Store.Set(req.Context(), state); err != nil {
		return response.NewError(err)
	}

	return response.NewRedirectWithStatus(
		m.Config.AuthCodeURL(state),
		response.HTTPStatusFound,
	)
}
