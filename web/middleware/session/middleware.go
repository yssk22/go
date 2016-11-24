package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/speedland/go/uuid"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xlog"
	"github.com/speedland/go/x/xtime"
)

// SessionLoggerKey is a logger key for the middleware
const SessionLoggerKey = "web.middleware.session"

type Middleware struct {
	Store SessionStore
	//
	// Configurations for cookie that keeps session id.
	//
	CookieName string
	MaxAge     time.Duration
	Domain     string
	HttpOnly   bool
	Path       string
}

// Default is a middleware instance with default value
var Default = &Middleware{
	Store:      NewMemorySessionStore(),
	CookieName: "speedland-go-session",
	MaxAge:     7 * 24 * time.Hour,
	Domain:     "localhost",
	HttpOnly:   true,
	Path:       "/",
}

// NewMiddleware returns the *Middleware with default configurations.
func NewMiddleware(store SessionStore) *Middleware {
	m := &Middleware{}
	*m = *Default
	m.Store = store
	return m
}

func (m *Middleware) Process(req *web.Request, next web.NextHandler) *response.Response {
	var logger = xlog.WithContext(req.Context()).WithKey(SessionLoggerKey)
	session, err := m.prepareSession(req)
	if err != nil {
		logger.Errorf("Failed to prepare sessoin(%v), fallback to create a new session", err)
		session = NewSession()
	}
	resp := next(req.WithContext(
		NewContext(req.Context(), session),
	))
	cookie, err := m.storeSession(req, session)
	if err != nil {
		logger.Errorf("Failed to store sessoin: %v", err)
		return resp
	}
	if cookie != nil {
		resp.Header.Set("X-SPEEDLAND-SESSION-ID", session.ID.String())
		resp.SetCookie(cookie, req.Option.HMACKey)
		if !session.fromStore {
			logger.Infof("Session initialized: %s", session.ID)
		}
	}
	return resp

}

func (m *Middleware) prepareSession(req *web.Request) (*Session, error) {
	strSessionID := req.Cookies.GetStringOr(m.CookieName, "")
	if strSessionID == "" {
		return NewSession(), nil
	}
	sessionID, ok := uuid.FromString(strSessionID)
	if !ok {
		return nil, fmt.Errorf("invalid session id: %q", strSessionID)
	}
	var logger = xlog.WithContext(req.Context()).WithKey(SessionLoggerKey)
	session, err := m.Store.Get(req.Context(), sessionID)
	logger.Debug(func(fmt *xlog.Printer) {
		fmt.Printf("Getting a session with %s\n", sessionID)
		fmt.Printf("From %v\n", m.Store)
		if session != nil {
			if buff, err := json.MarshalIndent(session, "", "\t"); err == nil {
				fmt.Println(string(buff))
			}
		}
	})
	if err != nil {
		return nil, fmt.Errorf("SessionStore.Get: %v", err)
	}
	session.fromStore = true
	if session.fromStore && session.IsExpired(m.MaxAge) {
		logger.Debugf("Session %s is expired, deleting", session.ID)
		m.Store.Del(req.Context(), session)
		return NewSession(), nil
	}
	return session, nil
}

func (m *Middleware) storeSession(req *web.Request, session *Session) (*http.Cookie, error) {
	var err error
	if session.fromStore ||
		xtime.Now().After(session.Timestamp.Add(m.MaxAge/2)) ||
		len(session.Data) > 0 {
		session.Timestamp = xtime.Now()
		err = m.Store.Set(req.Context(), session)
		return &http.Cookie{
			Name:     m.CookieName,
			Value:    session.ID.String(),
			Domain:   m.Domain,
			HttpOnly: m.HttpOnly,
			MaxAge:   int(m.MaxAge / time.Second),
			Path:     m.Path,
		}, err
	}
	return nil, nil
}
