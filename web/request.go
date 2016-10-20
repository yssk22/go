package web

import (
	"net/http"

	"github.com/speedland/go/keyvalue"
)

// Request is a http request.
type Request struct {
	*http.Request
	query  *keyvalue.GetProxy
	params *keyvalue.GetProxy
	form   *keyvalue.GetProxy
}

// Params returns an accessor to request parameters.
func (r *Request) Params() *keyvalue.GetProxy {
	return r.params
}

// Query returns an accessor to query parameters.
func (r *Request) Query() *keyvalue.GetProxy {
	return r.query
}

// Form returns an accessor to form values.
func (r *Request) Form() *keyvalue.GetProxy {
	return r.query
}
