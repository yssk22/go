package service

import (
	"net/http"

	"context"

	"google.golang.org/appengine/urlfetch"
)

// NewHTTPClient returns *http.Client under the gae service context
func NewHTTPClient(ctx context.Context) (c *http.Client) {
	c = http.DefaultClient
	s := FromContext(ctx)
	if s != nil {
		c = s.Config.NewHTTPClient(ctx)
		return
	}
	defer func() {
		recover()
		c = http.DefaultClient
	}()
	c = urlfetch.Client(ctx)
	return
}
