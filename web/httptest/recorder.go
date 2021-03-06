package httptest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/yssk22/go/x/xnet/xhttp/xhttptest"
)

// RequestFactory is an interface to create a new request.
type RequestFactory interface {
	NewRequest(string, string, io.Reader) (*http.Request, error)
}

// RequestFactoryFunc is a wrapper to NewRequest function to be RequestFactory
type RequestFactoryFunc func(string, string, io.Reader) (*http.Request, error)

// NewRequest implements RequestFactory#NewRequest
func (f RequestFactoryFunc) NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	return f(method, path, body)
}

// DefaultRequestFactory is a RequestFactory on top of http.NewRequest
var DefaultRequestFactory = RequestFactoryFunc(http.NewRequest)

// Recorder is a test recorder.
type Recorder struct {
	handler        http.Handler
	requestFactory RequestFactory
	Cookies        []*http.Cookie
}

// NewRecorder returns a new *Recorder
func NewRecorder(handler http.Handler) *Recorder {
	return &Recorder{
		handler:        handler,
		requestFactory: DefaultRequestFactory,
	}
}

// NewRecorderWithFactory returns a new *Recorder with the custom factory configuration
func NewRecorderWithFactory(handler http.Handler, factory RequestFactory) *Recorder {
	return &Recorder{
		handler:        handler,
		requestFactory: factory,
	}
}

// TestGet make a test GET request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestGet(path string, options ...RequestOption) *ResponseRecorder {
	w := httptest.NewRecorder()
	req := r.NewRequest("GET", path, nil, options...)
	r.handler.ServeHTTP(w, req)
	r.Cookies, _ = xhttptest.GetCookies(w)
	return &ResponseRecorder{w, req}
}

// TestPost make a test POST request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestPost(path string, v interface{}, options ...RequestOption) *ResponseRecorder {
	w := httptest.NewRecorder()
	req := r.NewRequest("POST", path, v, options...)
	r.handler.ServeHTTP(w, req)
	r.Cookies, _ = xhttptest.GetCookies(w)
	return &ResponseRecorder{w, req}
}

// TestPut make a test PUT request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestPut(path string, v interface{}, options ...RequestOption) *ResponseRecorder {
	w := httptest.NewRecorder()
	req := r.NewRequest("PUT", path, v, options...)
	r.handler.ServeHTTP(w, req)
	r.Cookies, _ = xhttptest.GetCookies(w)
	return &ResponseRecorder{w, req}
}

// TestDelete make a test DELETE request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestDelete(path string, options ...RequestOption) *ResponseRecorder {
	w := httptest.NewRecorder()
	req := r.NewRequest("DELETE", path, nil, options...)
	r.handler.ServeHTTP(w, req)
	r.Cookies, _ = xhttptest.GetCookies(w)
	return &ResponseRecorder{w, req}
}

// NewRequest returns http.Request with request body given by `v`
func (r *Recorder) NewRequest(method, path string, v interface{}, options ...RequestOption) *http.Request {
	var err error
	var req *http.Request
	var factory = r.requestFactory
	if v == nil {
		req, err = factory.NewRequest(method, path, nil)
	} else {
		switch v.(type) {
		case url.Values:
			req, err = factory.NewRequest(method, path, strings.NewReader(v.(url.Values).Encode()))
			if req != nil {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
		case io.Reader:
			req, err = factory.NewRequest(method, path, v.(io.Reader))
		default:
			var buff []byte
			buff, err = json.Marshal(v)
			if err != nil {
				panic(fmt.Errorf("Could not marshal the request body : %v (must be url.Values, io.Reader, or json marhslable.)", v))
			}
			req, err = factory.NewRequest(method, path, bytes.NewReader(buff))
		}
	}
	if err != nil {
		panic(fmt.Errorf("Could not prepare a request: %v", err))
	}
	for _, opts := range options {
		req = opts(req)
	}
	for _, c := range r.Cookies {
		req.AddCookie(c)
	}
	return req
}

// TestRequest make a test request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestRequest(req *http.Request) *ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, req)
	r.Cookies, _ = xhttptest.GetCookies(w)
	return &ResponseRecorder{w, req}
}

// RequestOption is a option function to configure a request
type RequestOption func(req *http.Request) *http.Request

// Header set the request header for Test functions
func Header(key string, value string) RequestOption {
	return func(req *http.Request) *http.Request {
		req.Header.Set(key, value)
		return req
	}
}
