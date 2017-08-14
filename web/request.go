package web

import (
	"fmt"
	"net/http"

	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/uuid"
	"github.com/speedland/go/x/xcontext"
	"github.com/speedland/go/x/xnet/xhttp"

	"context"
)

const defaultMaxMemory = 32 << 20 // 32 MB

// Request is a wrapper for net/http.Request
// The original `*net/http.Request` functions and fields are embedded in struct and provides
// some utility functions (especially to support context.Context)
type Request struct {
	*http.Request

	// common request scoped values
	ID      uuid.UUID
	Params  *keyvalue.GetProxy
	Query   *keyvalue.GetProxy
	Form    *keyvalue.GetProxy
	Cookies *keyvalue.GetProxy

	Option *Option
}

var requestContextKey = xcontext.NewKey("request")

// FromContext returns a *Request associated with the context.
func FromContext(ctx context.Context) *Request {
	req, ok := ctx.Value(requestContextKey).(*Request)
	if ok {
		return req
	}
	return nil
}

// NewRequest returns a new *Request
func NewRequest(r *http.Request, option *Option) *Request {
	if option == nil {
		option = DefaultOption
	}
	query := r.URL.Query()
	cookies := make(map[interface{}]*http.Cookie)
	for _, cc := range r.Cookies() {
		c, err := xhttp.UnsignCookie(cc, option.HMACKey)
		if err == nil {
			cookies[c.Name] = c
		}
	}
	if option.InitContext != nil {
		ctx := option.InitContext(r)
		r = r.WithContext(ctx)
	}
	req := &Request{
		Request: r,
		ID:      uuid.New(),
		Query:   keyvalue.NewQueryProxy(query),
		Form: keyvalue.GetterStringKeyFunc(func(key string) (interface{}, error) {
			if r.Form == nil {
				r.ParseMultipartForm(defaultMaxMemory)
			}
			if vs := r.Form[key]; len(vs) > 0 {
				return vs[0], nil
			}
			return nil, keyvalue.KeyError(key)
		}).Proxy(),
		Cookies: keyvalue.GetterStringKeyFunc(func(key string) (interface{}, error) {
			v, ok := cookies[key]
			if !ok {
				return nil, keyvalue.KeyError(key)
			}
			return v.Value, nil
		}).Proxy(),
		Option: option,
	}
	return req.WithValue(requestContextKey, req)
}

// WithContext returns a shallow copy of r with its context changed to ctx. The provided ctx must be non-nil.
func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("ctx must not be nil")
	}
	rr := new(Request)
	*rr = *r
	rr.Request = rr.Request.WithContext(ctx)
	return rr
}

// WithValue sets the request-scoped value with the in-flight http request and return a shallow copied request.
// This is shorthand for `req.WithContext(context.WithValue(req.Context(), key, value))`
func (r *Request) WithValue(key interface{}, value interface{}) *Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}

// Get implements keyvalue.Getter to enable keyvalue.GetProxy for context values.
func (r *Request) Get(key interface{}) (interface{}, error) {
	v := r.Context().Value(key)
	if v == nil {
		return nil, keyvalue.KeyError(fmt.Sprintf("%s", key))
	}
	return v, nil
}
