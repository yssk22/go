package retry

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPTransport is a type for http.Transport with retry
type HTTPTransport struct {
	Base    http.RoundTripper
	Checker HTTPChecker
	Backoff HTTPBackoff
}

// HTTPBackoff is an http version of Backoff
type HTTPBackoff interface {
	Calc(int, *http.Request, *http.Response, error) time.Duration
}

// HTTPChecker is an http version of Checker
type HTTPChecker interface {
	NeedRetry(int, *http.Request, *http.Response, error) bool
}

// NewHTTPTransport returns a new http.RoundTripper instance for given checker and backoff configurations on top of base http.RoundTripper
func NewHTTPTransport(base http.RoundTripper, checker HTTPChecker, backoff HTTPBackoff) http.RoundTripper {
	return &HTTPTransport{
		Base:    base,
		Checker: checker,
		Backoff: backoff,
	}
}

// RoundTrip implements http.Transport#RoundTrip
func (t *HTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var attempt int
	// request body should be buffered on memory since it is consumed by RoundTrip.
	reqBody, err := t.bufferBody(req)
	if err != nil {
		return nil, fmt.Errorf("could not allocate request body for *retry.HTTPTransport: %v", err)
	}
	for {
		if reqBody != nil {
			if _, serr := reqBody.Seek(0, 0); serr != nil {
				return nil, fmt.Errorf("could not reset request body for *retry.HTTPTransport: %v", err)
			}
			req.Body = ioutil.NopCloser(reqBody)
		}
		res, err := t.Base.RoundTrip(req)
		attempt++
		if !t.Checker.NeedRetry(attempt, req, res, err) {
			return res, err
		}
		// discard body content for retry.
		if res != nil && res.Body != nil {
			io.Copy(ioutil.Discard, res.Body)
			res.Body.Close()
		}
		ticker := time.NewTicker(t.Backoff.Calc(attempt, req, res, err))
		select {
		case <-req.Cancel:
			return nil, fmt.Errorf("request canceled")
		case <-ticker.C:
			break
		}
	}
}

func (t *HTTPTransport) bufferBody(req *http.Request) (*bytes.Reader, error) {
	if req.Body == nil {
		return nil, nil
	}
	defer req.Body.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, req.Body); err != nil {
		return nil, err
	}
	req.Body.Close()
	return bytes.NewReader(buf.Bytes()), nil
}

// HTTPConstBackoff is a http version of ConstBackoff
func HTTPConstBackoff(t time.Duration) HTTPBackoff {
	return &httpConstBackoff{
		interval: t,
	}
}

type httpConstBackoff struct {
	interval time.Duration
}

// Calc implements HTTPBackoff#Calc()
func (b *httpConstBackoff) Calc(int, *http.Request, *http.Response, error) time.Duration {
	return b.interval
}

// HTTPAnd is a AND combination of multiple HTTPChecker instances.
func HTTPAnd(checkers ...HTTPChecker) HTTPChecker {
	return &httpAnd{
		checkers: checkers,
	}
}

type httpAnd struct {
	checkers []HTTPChecker
}

func (and *httpAnd) NeedRetry(attempt int, req *http.Request, resp *http.Response, err error) bool {
	var needRetry = true
	for _, c := range and.checkers {
		needRetry = needRetry && c.NeedRetry(attempt, req, resp, err)
		if !needRetry {
			return false
		}
	}
	return true
}

// HTTPMaxRetry is an http version of MaxRetry
func HTTPMaxRetry(max int) HTTPChecker {
	return &httpMaxRetry{
		max: max,
	}
}

type httpMaxRetry struct {
	max int
}

func (mr *httpMaxRetry) NeedRetry(attempt int, req *http.Request, resp *http.Response, err error) bool {
	return attempt < mr.max
}

// HTTPServerErrorChecker returns a HTTPChecker that needs retries when http status code >= 500
func HTTPServerErrorChecker() HTTPChecker {
	return HTTPResponseChecker(
		func(resp *http.Response) bool {
			return resp.StatusCode >= 500
		},
	)
}

// HTTPResponseChecker returns a HTTPChecker that checks *http.Response for retries
func HTTPResponseChecker(f func(resp *http.Response) bool) HTTPChecker {
	return &httpResponseChecker{
		F: f,
	}
}

type httpResponseChecker struct {
	F func(resp *http.Response) bool
}

func (c *httpResponseChecker) NeedRetry(attempt int, req *http.Request, resp *http.Response, err error) bool {
	if err != nil || resp == nil {
		return true
	}
	return c.F(resp)
}
