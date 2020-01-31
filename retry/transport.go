package retry

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// httpTransport is a type for http.Transport with retry
type httpTransport struct {
	Base    http.RoundTripper
	Cond    HTTPRetryCond
	Backoff HTTPBackoff
}

// HTTPBackoff is an http version of Backoff
type HTTPBackoff interface {
	Calc(int, *http.Request, *http.Response, error) time.Duration
}

// HTTPRetryCond is an http version of Checker
type HTTPRetryCond interface {
	NeedRetry(int, *http.Request, *http.Response, error) bool
}

// NewHTTPTransport returns a new http.RoundTripper instance for given checker and backoff configurations on top of base http.RoundTripper
func NewHTTPTransport(base http.RoundTripper, cond HTTPRetryCond, backoff HTTPBackoff) http.RoundTripper {
	return &httpTransport{
		Base:    base,
		Cond:    cond,
		Backoff: backoff,
	}
}

// RoundTrip implements http.Transport#RoundTrip
func (t *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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
		if !t.Cond.NeedRetry(attempt, req, res, err) {
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

func (t *httpTransport) bufferBody(req *http.Request) (*bytes.Reader, error) {
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

// HTTPAnd is a AND combination of multiple HTTPRetryCond instances.
func HTTPAnd(checkers ...HTTPRetryCond) HTTPRetryCond {
	return &httpAnd{
		checkers: checkers,
	}
}

type httpAnd struct {
	checkers []HTTPRetryCond
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

// HTTPRetryUntil sets http request retries until max count
func HTTPRetryUntil(max int) HTTPRetryCond {
	return &httpRetryUntil{
		max: max,
	}
}

type httpRetryUntil struct {
	max int
}

func (mr *httpRetryUntil) NeedRetry(attempt int, req *http.Request, resp *http.Response, err error) bool {
	return attempt < mr.max
}

// HTTPRetryOnServerError returns a HTTPRetryCond that needs retries when http status code >= 500
func HTTPRetryOnServerError() HTTPRetryCond {
	return HTTPRetryIf(
		func(resp *http.Response) bool {
			return resp.StatusCode >= 500
		},
	)
}

// HTTPRetryIf returns a HTTPRetryCond that checks *http.Response for retries
func HTTPRetryIf(f func(resp *http.Response) bool) HTTPRetryCond {
	return &httpRetryIf{
		F: f,
	}
}

type httpRetryIf struct {
	F func(resp *http.Response) bool
}

func (c *httpRetryIf) NeedRetry(attempt int, req *http.Request, resp *http.Response, err error) bool {
	if err != nil || resp == nil {
		return true
	}
	return c.F(resp)
}
