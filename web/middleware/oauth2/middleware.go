package oauth2

import (
	"fmt"

	"github.com/speedland/go/uuid"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const OAuth2LoggerKey = "web.middleware.oauth2"

type Middleware struct {
	AuthPath        string // path that redirects to oauth2 provider
	CallbackPath    string
	AuthCodeOptions []oauth2.AuthCodeOption
	Store           StateStore
	Config          Config
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
	if m.CallbackPath == req.URL.EscapedPath() {
		return m.handleCallbackPath(req, next)
	}
	return next(req)
}

func (m *Middleware) handleAuthPath(req *web.Request) *response.Response {
	state := uuid.New().String()
	if err := m.Store.Set(req.Context(), state); err != nil {
		return response.NewError(err)
	}

	return response.NewRedirectWithStatus(
		m.Config.AuthCodeURL(state, m.AuthCodeOptions...),
		response.HTTPStatusFound,
	)
}

func (m *Middleware) handleCallbackPath(req *web.Request, next web.NextHandler) *response.Response {
	code := req.Form.GetStringOr("code", "")
	if code == "" {
		return response.NewErrorWithStatus(
			fmt.Errorf("code is required"),
			response.HTTPStatusBadRequest,
		)
	}
	state := req.Form.GetStringOr("state", "")
	if state == "" {
		return response.NewErrorWithStatus(
			fmt.Errorf("state is required"),
			response.HTTPStatusBadRequest,
		)
	}
	storedState, err := m.Store.Get(req.Context())
	if err != nil {
		return response.NewError(
			fmt.Errorf("validation failure: %v", err),
		)
	}
	if state != storedState {
		return response.NewErrorWithStatus(
			fmt.Errorf("invalid state"),
			response.HTTPStatusBadRequest,
		)
	}
	token, err := m.Config.Exchange(req.Context(), code)
	if err != nil {
		return response.NewErrorWithStatus(
			fmt.Errorf("failed to exchange the token: %v", err),
			response.HTTPStatusBadRequest,
		)
	}
	// TODO: token handling.
	return next(req.WithContext(context.WithValue(req.Context(), tokenContextKey, token)))
}
