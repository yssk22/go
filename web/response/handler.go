package response

import (
	"context"
	"io"
	"net/http"
)

type _handler struct {
	request *http.Request
	handler http.Handler
}

func (h *_handler) Render(ctx context.Context, w io.Writer) {
	h.handler.ServeHTTP(w.(http.ResponseWriter), h.request)
}

// NewHandler returns a new *Response that calls h.ServeHTTP
func NewHandler(ctx context.Context, h http.Handler, req *http.Request) *Response {
	return NewResponseWithStatus(ctx, &_handler{request: req, handler: h}, 200)
}
