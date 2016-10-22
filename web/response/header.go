package response

import (
	"net/http"

	"golang.org/x/net/context"
)

// Header is a Response implementation that writes response header
type Header struct {
	Code   HTTPStatus
	Fields http.Header
}

// NewHeader returns a new Header
func NewHeader() *Header {
	return &Header{
		Code: HTTPStatusOK,
		Fields: http.Header(
			make(map[string][]string),
		),
	}
}

// ContentType sets the `Content-Type` header.
func (h *Header) ContentType(t string) *Header {
	h.Fields.Set("content-type", t)
	return h
}

// Render renders the Fields to `w`
func (h *Header) Render(ctx context.Context, w http.ResponseWriter) {
	header := w.Header()
	for key, v := range h.Fields {
		for _, vv := range v {
			header.Add(key, vv)
		}
	}
	// merge headers in ctx.
	if ctxHeader, ok := ctx.Value(responseHeaderKey).(http.Header); ok {
		for k, v := range ctxHeader {
			for _, vv := range v {
				header.Add(k, vv)
			}
		}
	}
	w.WriteHeader(int(h.Code))
}

var responseHeaderKey = &contextKey{"header"}

// SetHeader sets the header field on the current context.
// You can set the specific respnse header in your handler:
//
//     req.WithContext(
//        response.SetHeader(req.Context(), key, value),
//     )
//
//  Note:
//
//     This interface may change since it's a bit strange to store response header into
//     request context.
//
func SetHeader(ctx context.Context, key string, value string) context.Context {
	header, ok := ctx.Value(responseHeaderKey).(http.Header)
	if !ok {
		header = http.Header(make(map[string][]string))
		ctx = context.WithValue(ctx, responseHeaderKey, header)
	}
	header.Set(key, value)
	return ctx
}

// AddHeader is like SetHeader but adding instead of setting.
func AddHeader(ctx context.Context, key string, value string) context.Context {
	header, ok := ctx.Value(responseHeaderKey).(http.Header)
	if !ok {
		header = http.Header(make(map[string][]string))
	}
	header.Add(key, value)
	return ctx
}
