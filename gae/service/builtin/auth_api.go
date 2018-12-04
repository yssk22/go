package builtin

import (
	"context"

	"github.com/yssk22/go/web/api"
	"github.com/yssk22/go/gae/service"
	"github.com/yssk22/go/gae/service/auth"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xlog"
	"google.golang.org/appengine"
)

// @api path=/auth/me.json
func getMe(ctx context.Context) (*auth.Auth, error) {
	s := service.FromContext(ctx)
	ctx, err := appengine.Namespace(ctx, s.APIConfig.AuthNamespace)
	xerrors.MustNil(err)
	return auth.GetCurrent(ctx)
}

type fbLoginParams struct {
	AccessToken string `json:"access_token"`
}

// @api path=/auth/login/facebook/
func createFacebookAuth(ctx context.Context, params *fbLoginParams) (*auth.Auth, error) {
	s := service.FromContext(ctx)
	ctx, err := appengine.Namespace(ctx, s.APIConfig.AuthNamespace)
	xerrors.MustNil(err)
	ctx, logger := xlog.WithContext(ctx, "")
	a, err := auth.Facebook(ctx, s.Config.NewHTTPClient(ctx), params.AccessToken)
	if err != nil {
		logger.Infof("failed to authenticate facebook with %q: %v", params.AccessToken, err)
		return nil, api.BadRequest
	}
	xerrors.MustNil(auth.SetCurrent(ctx, a))
	return a, nil
}

// @api method=GET path=/auth/logout/
func processLoggedOutForDev(ctx context.Context) (bool, error) {
	if !appengine.IsDevAppServer() {
		// available only in DevAppServer
		return false, api.BadRequest
	}
	auth.DeleteCurrent(ctx)
	return true, nil
}

// @api method=POST path=/auth/logout/
func processLoggedOut(ctx context.Context) (bool, error) {
	auth.DeleteCurrent(ctx)
	return true, nil
}
