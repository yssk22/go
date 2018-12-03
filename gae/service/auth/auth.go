package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yssk22/go/services/facebook"
	"github.com/yssk22/go/uuid"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xtime"
)

// AuthType is an enum for authentication types.
// @ent
type AuthType int

// AuthType values
const (
	AuthTypeNone AuthType = iota
	AuthTypeEmail
	AuthTypeFacebook
	AuthTypeTwitter
	AuthTypeMessenger
	AuthTypeGoogle
)

// Auth is a primary type to represent a user
// @datastore
type Auth struct {
	ID          string    `json:"id" ent:"id"`
	FacebookID  string    `json:"facebook_id"`
	MessengerID string    `json:"messenger_id"`
	TwitterID   string    `json:"twitter_id"`
	Email       string    `json:"email"`
	AuthType    AuthType  `json:"auth_type"`
	IsAdmin     bool      `json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
	LastLoginAt time.Time `json:"last_login_at"`
	UpdatedAt   time.Time `json:"updated_at" ent:"timestamp"`
}

// Guest is a unauthorized value.
var Guest = &Auth{
	ID:          "guest",
	FacebookID:  "guest",
	MessengerID: "guest",
	TwitterID:   "guest",
	Email:       "",
	AuthType:    AuthTypeNone,
	CreatedAt:   time.Time{},
	LastLoginAt: time.Time{},
	UpdatedAt:   time.Time{},
}

// IsGuest returns a *Auth is guest auth or not.
func (a *Auth) IsGuest() bool {
	return a.ID == Guest.ID
}

var (
	fbAuthContextKey = struct{}{}
)

// Facebook process authorization by the given access token
func Facebook(ctx context.Context, client *http.Client, token string) (*Auth, error) {
	user, err := GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("session error: %v", err)
	}
	fb := facebook.NewClient(client, token)
	me, err := fb.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user, err = findExistingAuth(ctx, AuthTypeFacebook, me.ID)
		if err != nil {
			return nil, err
		}
	}
	user.FacebookID = me.ID
	user.LastLoginAt = xtime.Now()
	if _, err = NewAuthKind().Put(ctx, user); err != nil {
		return nil, fmt.Errorf("datastore error: %v", err)
	}
	if err = SetCurrent(ctx, user); err != nil {
		return nil, fmt.Errorf("session error: %v", err)
	}
	return user, nil
}

var (
	messengerAuthContextKey = struct{}{}
)

// Messenger process authorization by the given access token
func Messenger(ctx context.Context, messengerID string) (*Auth, error) {
	user, err := GetCurrent(ctx)
	if err != nil {
		return nil, xerrors.Wrap(err, "session error to get the current session")
	}
	if user == nil {
		user, err = findExistingAuth(ctx, AuthTypeMessenger, messengerID)
		if err != nil {
			return nil, err
		}
	}
	user.MessengerID = messengerID
	user.LastLoginAt = xtime.Now()
	if _, err = NewAuthKind().Put(ctx, user); err != nil {
		return nil, fmt.Errorf("datastore error: %v", err)
	}
	if err = SetCurrent(ctx, user); err != nil {
		return nil, fmt.Errorf("session error: %v", err)
	}
	return user, nil
}

func findExistingAuth(ctx context.Context, at AuthType, id string) (*Auth, error) {
	q := NewAuthQuery().EqAuthType(at)
	switch at {
	case AuthTypeFacebook:
		q = q.EqFacebookID(id)
	default:

	}
	_, values, err := q.Limit(1).GetAll(ctx)
	if err != nil {
		return nil, xerrors.Wrap(err, "could not query the auth records")
	}
	if len(values) >= 1 {
		return &values[0], nil
	}
	return &Auth{
		ID:        uuid.New().String(),
		AuthType:  at,
		CreatedAt: xtime.Now(),
	}, nil
}
