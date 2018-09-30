package auth

import (
	"fmt"
	"net/http"
	"time"

	"context"

	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/lazy"
	"github.com/yssk22/go/services/facebook"
	"github.com/yssk22/go/uuid"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xtime"
)

// AuthType is an enum for authentication types.
//go:generate enum -type=AuthType
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
//go:generate ent -type=Auth
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
	fbAuthQuery      = NewAuthQuery().Eq("FacebookID", lazy.Func(func(ctx context.Context) (interface{}, error) {
		val := ctx.Value(fbAuthContextKey)
		if me, ok := val.(*facebook.Me); ok {
			return me.ID, nil
		}
		return nil, keyvalue.KeyError("fbAuthContextKey")
	})).Limit(lazy.New(1))
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
		user, err = findExistingAuth(context.WithValue(ctx, fbAuthContextKey, me), fbAuthQuery, AuthTypeFacebook)
		if err != nil {
			return nil, err
		}
	}
	user.FacebookID = me.ID
	user.LastLoginAt = xtime.Now()
	if _, err = DefaultAuthKind.Put(ctx, user); err != nil {
		return nil, fmt.Errorf("datastore error: %v", err)
	}
	if err = SetCurrent(ctx, user); err != nil {
		return nil, fmt.Errorf("session error: %v", err)
	}
	return user, nil
}

var (
	messengerAuthContextKey = struct{}{}
	messengerAuthQuery      = NewAuthQuery().Eq("MessengerID", lazy.Func(func(ctx context.Context) (interface{}, error) {
		val := ctx.Value(messengerAuthContextKey)
		if id, ok := val.(string); ok {
			return id, nil
		}
		return nil, keyvalue.KeyError("messengerAuthContextKey")
	})).Limit(lazy.New(1))
)

// Messenger process authorization by the given access token
func Messenger(ctx context.Context, messengerID string) (*Auth, error) {
	user, err := GetCurrent(ctx)
	if err != nil {
		return nil, xerrors.Wrap(err, "session error to get the current session")
	}
	if user == nil {
		user, err = findExistingAuth(context.WithValue(ctx, messengerAuthContextKey, messengerID), messengerAuthQuery, AuthTypeMessenger)
		if err != nil {
			return nil, err
		}
	}
	user.MessengerID = messengerID
	user.LastLoginAt = xtime.Now()
	if _, err = DefaultAuthKind.Put(ctx, user); err != nil {
		return nil, xerrors.Wrap(err, "datastore error")
	}
	if err = SetCurrent(ctx, user); err != nil {
		return nil, xerrors.Wrap(err, "session error to set the current session")
	}
	return user, nil
}

func findExistingAuth(ctx context.Context, query *AuthQuery, at AuthType) (*Auth, error) {
	values, err := query.GetAllValues(ctx)
	if err != nil {
		return nil, fmt.Errorf("datastore error: %v", err)
	}
	if len(values) >= 1 {
		return values[0], nil
	}
	return &Auth{
		ID:        uuid.New().String(),
		AuthType:  at,
		CreatedAt: xtime.Now(),
	}, nil
}
