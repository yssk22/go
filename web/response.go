package web

import (
	"net/http"
)

// Response is an interface for htp response
// Implementations are available in github.com/speedland/go/web/response package
type Response interface {
	// Render should set the header and write the contents.
	// The return value should tell the end of content or not.
	Render(w http.ResponseWriter) bool
}
