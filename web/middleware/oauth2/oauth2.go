// Package oauth2 provides oauth2 middleware
package oauth2

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// StateStore is an interface for session storage.
type StateStore interface {
	Set(context.Context, string) error
	Validate(context.Context, string) (bool, error)
}

// Config is an interface for OAuth2 configuration. Use "golang.org/x/oauth2.Config"
type Config interface {
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
}
