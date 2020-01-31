package response

import (
	"bytes"
	"io"
	"net/http"

	"github.com/yssk22/go/x/xcrypto/xhmac"
	"github.com/yssk22/go/x/xnet/xhttp"

	"context"
)

// Response represents http response.
type Response struct {
	Status  HTTPStatus
	Header  http.Header
	Cookies []*http.Cookie
	Body    Body
	ctx     context.Context // context when the response is created
}

// NewResponse retuurns a *Response to write body content
func NewResponse(ctx context.Context, body Body) *Response {
	return &Response{
		Status:  HTTPStatusOK,
		Header:  http.Header(make(map[string][]string)),
		Cookies: []*http.Cookie{},
		Body:    body,
		ctx:     ctx,
	}
}

// NewResponseWithStatus retuurns a *Response to write body content with a custom status code.
func NewResponseWithStatus(ctx context.Context, body Body, status HTTPStatus) *Response {
	return &Response{
		Status: status,
		Header: http.Header(make(map[string][]string)),
		Body:   body,
		ctx:    ctx,
	}
}

// SetCookie add a cookie on the response
func (r *Response) SetCookie(c *http.Cookie, hmac *xhmac.Base64) {
	r.Cookies = append(r.Cookies, xhttp.SignCookie(c, hmac))
}

// Render renders whole http contnet
func (r *Response) Render(w http.ResponseWriter) {
	wh := w.Header()
	for k, v := range r.Header {
		for _, vv := range v {
			wh.Add(k, vv)
		}
	}
	for _, c := range r.Cookies {
		http.SetCookie(w, c)
	}
	w.WriteHeader(int(r.Status))
	r.Body.Render(r.ctx, w)
}

// Content returns the rendered result of the response body
func (r *Response) Content() string {
	var buff bytes.Buffer
	r.Body.Render(context.Background(), &buff)
	return buff.String()
}

// Context returns the context when the response is created
func (r *Response) Context() context.Context {
	return r.ctx
}

// Body is an interface to write response
type Body interface {
	Render(ctx context.Context, w io.Writer)
}

type noContent struct{}

func (r noContent) Render(ctx context.Context, w io.Writer) {
}

var NoContent = &noContent{}
