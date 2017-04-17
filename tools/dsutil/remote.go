package dsutil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/appengine"
	"google.golang.org/appengine/remote_api"
)

var remoteCtxScopes = []string{
	"https://www.googleapis.com/auth/appengine.apis",
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/cloud-platform",
}

var ErrNoRemoteAPIKeyIsConfigured = fmt.Errorf("no GOOGLE_APPLICATION_CREDENTIALS environment key is configured")

func GetRemoteContext(ctx context.Context, host string, namespace string, keyfile string) (context.Context, error) {
	var hc *http.Client
	var err error
	if keyfile != "" {
		jsonKey, err := ioutil.ReadFile(keyfile)
		if err != nil {
			return nil, err
		}
		cfg, err := google.JWTConfigFromJSON(jsonKey, remoteCtxScopes...)
		if err != nil {
			return nil, err
		}
		hc = oauth2.NewClient(ctx, cfg.TokenSource(ctx))
	} else {
		const defaultCredentialEnvKey = "GOOGLE_APPLICATION_CREDENTIALS"
		if os.Getenv(defaultCredentialEnvKey) != "" {
			hc, err = google.DefaultClient(oauth2.NoContext, remoteCtxScopes...)
		} else {
			return nil, ErrNoRemoteAPIKeyIsConfigured
		}
		if err != nil {
			return nil, err
		}
	}

	remoteCtx, err := remote_api.NewRemoteContext(host, hc)
	if err != nil {
		return nil, err
	}
	if namespace != "" {
		return appengine.Namespace(remoteCtx, namespace)
	}
	return remoteCtx, nil
}
