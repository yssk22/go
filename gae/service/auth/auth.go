package auth

import (
	"net/http"
	"time"

	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/lazy"
	"github.com/speedland/go/services/facebook"
	"github.com/speedland/go/uuid"
	"github.com/speedland/go/x/xtime"
	"golang.org/x/net/context"
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
)

// Auth is a primary type to represent a user
//go:generate ent -type=Auth
type Auth struct {
	ID          string    `json:"id" ent:"id"`
	FacebookID  string    `json:"facebook_id"`
	TwitterID   string    `json:"twitter_id"`
	Email       string    `json:"email"`
	AuthType    AuthType  `json:"auth_type"`
	CreatedAt   time.Time `json:"created_at"`
	LastLoginAt time.Time `json:"last_login_at"`
	UpdatedAt   time.Time `json:"updated_at" ent:"timestamp"`
}

// Guest is a unauthorized value.
var Guest = &Auth{
	ID:          "guest",
	FacebookID:  "guest",
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
	fb := facebook.NewClient(client, token)
	me, err := fb.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	var a *Auth
	var loginAt = xtime.Now()
	values, err := fbAuthQuery.GetAllValues(context.WithValue(ctx, fbAuthContextKey, me))
	if err != nil {
		return nil, err
	}
	if len(values) == 1 {
		a = values[0]
		a.LastLoginAt = loginAt
	} else {
		a = &Auth{
			ID:          uuid.New().String(),
			AuthType:    AuthTypeFacebook,
			FacebookID:  me.ID,
			LastLoginAt: loginAt,
			CreatedAt:   loginAt,
		}
	}
	if _, err = DefaultAuthKind.Put(ctx, a); err != nil {
		return nil, err
	}
	if err = SetCurrent(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}
