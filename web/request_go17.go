// +build go1.7

package web

import (
	"net/http"

	"golang.org/x/net/context"
)

func initContext(req *http.Request) context.Context {
	return req.Context()
}
