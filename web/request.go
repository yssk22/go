package web

import (
	"fmt"
	"net/http"

	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/uuid"

	"golang.org/x/net/context"
)

// Request is a wrapper for net/http.Request
// The original `*net/http.Request` functions and fields are embedded in struct and provides
// some utility functions (especially to support context.Context)
type Request struct {
	*http.Request
	ctx context.Context

	// common request scoped values
	ID     uuid.UUID
	Params *keyvalue.GetProxy
	Query  *keyvalue.GetProxy
	Form   *keyvalue.GetProxy
}

// NewRequest returns a new *Request
func NewRequest(r *http.Request) *Request {
	return &Request{
		Request: r,
		ctx:     initContext(r),
		ID:      uuid.New(),
		Query:   newURLValuesProxy(r.URL.Query()),
	}
}

// Context returns the context associated with request.
func (r *Request) Context() context.Context {
	return r.ctx
}

// WithContext returns a shallow copy of r with its context changed to ctx. The provided ctx must be non-nil.
func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("ctx must not be nil")
	}
	rr := new(Request)
	*rr = *r
	rr.ctx = ctx
	return rr
}

// WithValue sets the request-scoped value with the in-flight http request and return a shallow copied request.
// This is shorthand for `req.WithContext(context.WithValue(req.Context(), key, value))`
func (r *Request) WithValue(key interface{}, value interface{}) *Request {
	return r.WithContext(context.WithValue(r.ctx, key, value))
}

// Get implements keyvalue.Getter to enable keyvalue.GetProxy for context values.
func (r *Request) Get(key interface{}) (interface{}, error) {
	v := r.Context().Value(key)
	if v == nil {
		return nil, keyvalue.KeyError(fmt.Sprintf("%s", key))
	}
	return v, nil
}
