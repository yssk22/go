package xhttptest

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
)

// Stub enforce stub accesses for http.Client using URL mapping.
func Stub(mapping map[string]string, c *http.Client) *http.Client {
	rewriteMapping := make(map[string]*url.URL)
	for k, v := range mapping {
		var err error
		rewriteMapping[k], err = url.Parse(v)
		if err != nil {
			panic(fmt.Errorf("invalid stub url mapping: %s -> %s", k, v))
		}
	}
	c.Transport = &rewriter{
		mapping: rewriteMapping,
		base:    c.Transport,
	}
	return c
}

// rewriter is to rewrite a request not to access a remote resource
type rewriter struct {
	mapping map[string]*url.URL
	base    http.RoundTripper
}

// RoundTrip implements http.RoundTripper#RoundTrip
func (r *rewriter) RoundTrip(req *http.Request) (*http.Response, error) {
	u := req.URL.String()
	rewriteTo := r.mapping[u]
	if rewriteTo == nil {
		return nil, fmt.Errorf("forbitten by xhttptest.Stub")
	}
	req.URL = rewriteTo
	if r.base != nil {
		return r.base.RoundTrip(req)
	}
	return http.DefaultTransport.RoundTrip(req)
}

// StubFile is like a Stub but access local file resources
// instead of making actual HTTP requests
// The following fields are valid and others should not be used in your test.
//
//   - .Status
//   - .StatusCode
//   - .Header.Get("content-type")
//   - .ContentLength
//   - .Body
//
func StubFile(mapping map[string]string, c *http.Client) *http.Client {
	c.Transport = &fileStub{
		mapping: mapping,
	}
	return c
}

type fileStub struct {
	mapping map[string]string
}

// RoundTrip implements http.RoundTripper#RoundTrip
func (r *fileStub) RoundTrip(req *http.Request) (*http.Response, error) {
	u := req.URL.String()
	filePath := r.mapping[u]
	if filePath == "" {
		return nil, fmt.Errorf("forbitten by xhttptest.StubFile")
	}
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Status:     "OK",
		Close:      false,
		Header:     make(map[string][]string),
	}
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("stub file error: %v", err)
	}
	resp.ContentLength = stat.Size()
	resp.Header.Set("content-type", mime.TypeByExtension(path.Ext(filePath)))
	file, _ := os.Open(filePath)
	resp.Body = file
	return resp, nil
}
