// +build linux darwin freebsd netbsd openbsd dragonfly
// +build !appengine

package main

import (
	"golang.org/x/net/context"

	"github.com/speedland/go/services/google/photo"
	"github.com/speedland/go/x/xtime"
	"golang.org/x/oauth2"
)

const gphotoTokenType = "Bearer"
const gphotoAuthURL = "https://accounts.google.com/o/oauth2/auth"
const gphotoTokenURL = "https://accounts.google.com/o/oauth2/token"
const gphotoScopes = "https://picasaweb.google.com/data/"

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var DefaultConfig = &Config{}

func NewClient() *photo.Client {
	cfg := photo.NewOAuth2Config(DefaultConfig.ClientID, DefaultConfig.ClientSecret)
	token := &oauth2.Token{
		AccessToken:  DefaultConfig.AccessToken,
		RefreshToken: DefaultConfig.RefreshToken,
		TokenType:    "Bearer",
		// need Expiry otherwise we'll see "Invalid stateless token"
		Expiry: xtime.Now(),
	}
	return photo.NewClient(cfg.Client(context.Background(), token))
}
