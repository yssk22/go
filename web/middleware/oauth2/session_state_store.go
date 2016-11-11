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

// Get implements StateStore#Get
func (*SessionStateStore) Get(ctx context.Context) (string, error) {
	session := session.FromContext(ctx)
	if session == nil {
		return "", ErrSessionNotInitialized
	}
	var state string
	if err := session.Get(oauth2SessionStateKey, &state); err != nil {
		return "", fmt.Errorf("session(%q): %v", session.ID, err)
	}
	session.Del(oauth2SessionStateKey)
	return state, nil
}
