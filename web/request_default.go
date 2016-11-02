// +build !appengine
// +build !go1.7

package web

import (
	"net/http"

	"google.golang.org/appengine"

	"golang.org/x/net/context"
)

func initContext(req *http.Request) context.Context {
	var ctx = context.Background()
	// TODO: remove this if we find the clear solution.
	// use appengine.NewContext() here for `go test` to executes appengine related tests
	// and recover using context.Background() if context initialization failure due to non-appengine environment.
	func() {
		defer func() {
			recover()
		}()
		ctx = appengine.NewContext(req)
	}()
	return ctx
}
