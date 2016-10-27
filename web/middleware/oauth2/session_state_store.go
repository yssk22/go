package oauth2

import (
	"fmt"

	"github.com/speedland/go/web/middleware/session"
	"golang.org/x/net/context"
)

const oauth2SessionStateKey = "web.middleware.oauth2.state"

// SessionStateStore is an implementation of StateStore using session.
type SessionStateStore struct {
}

// ErrSessionNotInitialized is an error if the session middleware is not initialized.
var ErrSessionNotInitialized = fmt.Errorf("session is not initialized in the context")

// Set implements StateStore#Set
func (*SessionStateStore) Set(ctx context.Context, state string) error {
	session := session.FromContext(ctx)
	if session == nil {
		return ErrSessionNotInitialized
	}
	return session.Set(oauth2SessionStateKey, state)
}

// Set implements StateStore#Validate
func (*SessionStateStore) Validate(ctx context.Context, state string) (bool, error) {
	session := session.FromContext(ctx)
	if session == nil {
		return false, ErrSessionNotInitialized
	}
	v, err := session.Get(oauth2SessionStateKey)
	if err != nil {
		return false, fmt.Errorf("session(%q): %v", session.ID, err)
	}
	_, ok := v.(string)
	if !ok {
		return false, fmt.Errorf("OAuth2 state contains non-string: %v", v)
	}
	session.Del(oauth2SessionStateKey)
	return true, nil
}
