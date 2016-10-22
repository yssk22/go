package web

import (
	"net/http"

	"golang.org/x/net/context"
)

// Response is an interface to write response
type Response interface {
	Render(ctx context.Context, w http.ResponseWriter)
}
