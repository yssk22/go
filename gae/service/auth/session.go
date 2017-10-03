package auth

import (
	"fmt"

	"context"

	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/web/middleware/session"
	"google.golang.org/appengine/user"
)

const sessionKey = "current_user"

var ErrNoSession = fmt.Errorf("no session is available in the context")

// SetCurrent it to set the auth into the current session
func SetCurrent(ctx context.Context, a *Auth) error {
	s := session.FromContext(ctx)
	if s == nil {
		return ErrNoSession
	}
	return Set(s, a)
}

// GetCurrent is to get the auth fron the current session
func GetCurrent(ctx context.Context) (*Auth, error) {
	u := user.Current(ctx)
	if u != nil {
		return &Auth{
			ID:       fmt.Sprintf("gae_%s", u.ID),
			IsAdmin:  u.Admin,
			AuthType: AuthTypeGoogle,
		}, nil
	}
	s := session.FromContext(ctx)
	if s == nil {
		return nil, ErrNoSession
	}
	return Get(s)
}

// DeleteCurrent is to delete the auth fron the current session
func DeleteCurrent(ctx context.Context) error {
	s := session.FromContext(ctx)
	if s == nil {
		return ErrNoSession
	}
	return Delete(s)
}

// Set is to set the auth into the session.
func Set(s *session.Session, a *Auth) error {
	return s.Set(sessionKey, a)
}

// Get is to get the auth from the session
func Get(s *session.Session) (*Auth, error) {
	var a Auth
	if err := s.Get(sessionKey, &a); err != nil {
		if keyvalue.IsKeyError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

// Delete is to delete the auth from the session
func Delete(s *session.Session) error {
	return s.Del(sessionKey)
}
