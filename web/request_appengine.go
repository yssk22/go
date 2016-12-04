// +build appengine

package web

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func initContext(req *http.Request) context.Context {
	var ctx context.Context
	func() {
		defer func() {
			if x := recover(); x != nil {
				ctx = context.Background()
			}
		}()
		ctx = appengine.NewContext(req)
	}()
	return ctx
}
