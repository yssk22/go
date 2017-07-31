package response

import (
	"bytes"
	"io"
	"net/http"

	"github.com/speedland/go/x/xcrypto/xhmac"
	"github.com/speedland/go/x/xnet/xhttp"

	"context"
)

// Response represents http response.
type Response struct {
	Status  HTTPStatus
	Header  http.Header
	Cookies []*http.Cookie
	Body    Body
}

// NewResponse retuurns a *Response to write body content
func NewResponse(body Body) *Response {
	return &Response{
		Status:  HTTPStatusOK,
		Header:  http.Header(make(map[string][]string)),
		Cookies: []*http.Cookie{},
		Body:    body,
	}
}

// NewResponseWithStatus retuurns a *Response to write body content with a custom status code.
func NewResponseWithStatus(body Body, status HTTPStatus) *Response {
	return &Response{
		Status: status,
		Header: http.Header(make(map[string][]string)),
		Body:   body,
	}
}

// SetCookie add a cookie on the response
func (r *Response) SetCookie(c *http.Cookie, hmac *xhmac.Base64) {
	r.Cookies = append(r.Cookies, xhttp.SignCookie(c, hmac))
}

// Render renders whole http contnet
func (r *Response) Render(ctx context.Context, w http.ResponseWriter) {
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
	r.Body.Render(ctx, w)
}

// Content returns the rendered result of the response body
func (r *Response) Content() string {
	var buff bytes.Buffer
	r.Body.Render(context.Background(), &buff)
	return buff.String()
}

// Body is an interface to write response
type Body interface {
	Render(ctx context.Context, w io.Writer)
}

type noContent struct{}

func (r noContent) Render(ctx context.Context, w io.Writer) {
}

var NoContent = &noContent{}
