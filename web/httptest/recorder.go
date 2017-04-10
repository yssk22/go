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

	"github.com/speedland/go/x/xnet/xhttp/xhttptest"
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
func (r *Recorder) TestGet(path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, r.NewRequest("GET", path, nil))
	r.Cookies, _ = xhttptest.GetCookies(w)
	return w
}

// TestPost make a test POST request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestPost(path string, v interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, r.NewRequest("POST", path, v))
	r.Cookies, _ = xhttptest.GetCookies(w)
	return w
}

// TestPut make a test PUT request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestPut(path string, v interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, r.NewRequest("PUT", path, v))
	r.Cookies, _ = xhttptest.GetCookies(w)
	return w
}

// TestDelete make a test DELETE request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestDelete(path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, r.NewRequest("DELETE", path, nil))
	r.Cookies, _ = xhttptest.GetCookies(w)
	return w
}

// PrepareRequest prepares http.Request with request body
func (r *Recorder) NewRequest(method, path string, v interface{}) *http.Request {
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
			req, err = factory.NewRequest("POST", path, bytes.NewReader(buff))
		}
	}
	if err != nil {
		panic(fmt.Errorf("Could not prepare a request: %v", err))
	}
	for _, c := range r.Cookies {
		req.AddCookie(c)
	}
	return req
}

// TestRequest make a test request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestRequest(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, req)
	r.Cookies, _ = xhttptest.GetCookies(w)
	return w
}
