// +build appengine

package web

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func initContext(r *http.Request) context.Context {
	return appengine.NewContext(r.Request)
}
