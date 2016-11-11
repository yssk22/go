// Package oauth2 provides oauth2 middleware
package oauth2

import (
	"github.com/speedland/go/x/xcontext"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// StateStore is an interface for session storage.
type StateStore interface {
	Set(context.Context, string) error
	Get(context.Context) (string, error)
}

// Config is an interface for OAuth2 configuration. Use "golang.org/x/oauth2.Config"
type Config interface {
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
	Exchange(context.Context, string) (*oauth2.Token, error)
}

var tokenContextKey = xcontext.NewKey("token")

func FromContext(ctx context.Context) *oauth2.Token {
	token, ok := ctx.Value(tokenContextKey).(*oauth2.Token)
	if !ok {
		return nil
	}
	return token
}
