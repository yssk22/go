package service

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

// NewHTTPClient returns *http.Client under the gae service context
func NewHTTPClient(ctx context.Context) *http.Client {
	s := FromContext(ctx)
	if s == nil {
		return s.Config.NewHTTPClient(ctx)
	}
	return urlfetch.Client(ctx)
}
