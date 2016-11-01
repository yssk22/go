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
)

// Recorder is a test recorder.
type Recorder struct {
	handler http.Handler
}

// NewRouter returns *Router
func NewRecorder(handler http.Handler) *Recorder {
	return &Recorder{
		handler: handler,
	}
}

// TestGet make a test GET request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestGet(path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, NewRequest("GET", path, nil))
	return w
}

// TestPost make a test POST request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestPost(path string, v interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, NewRequest("POST", path, v))
	return w
}

// TestPut make a test PUT request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestPut(path string, v interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, NewRequest("PUT", path, v))
	return w
}

// TestDelete make a test DELETE request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestDelete(path string, v interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, NewRequest("DELETE", path, nil))
	return w
}

// PrepareRequest prepares http.Request with request body
func NewRequest(method, path string, v interface{}) *http.Request {
	var err error
	var req *http.Request
	if v == nil {
		req, err = http.NewRequest(method, path, nil)
	} else {
		switch v.(type) {
		case url.Values:
			req, err = http.NewRequest(method, path, strings.NewReader(v.(url.Values).Encode()))
			if req != nil {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
		case io.Reader:
			req, err = http.NewRequest(method, path, v.(io.Reader))
		default:
			var buff []byte
			buff, err = json.Marshal(v)
			if err != nil {
				panic(fmt.Errorf("Could not marshal the request body : %v (must be url.Values, io.Reader, or json marhslable.)", v))
			}
			req, err = http.NewRequest("POST", path, bytes.NewReader(buff))
		}
	}
	if err != nil {
		panic(fmt.Errorf("Could not prepare a request: %v", err))
	}
	return req
}

// TestRequest make a test request to the router and returns the response as *http.ResponseRecorder
func (r *Recorder) TestRequest(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, req)
	return w
}
