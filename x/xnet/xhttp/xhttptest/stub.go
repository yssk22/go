package xhttptest

import (
	"fmt"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/speedland/go/x/xerrors"
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

// StubServer is an http server that serve the request as stub.
type StubServer struct {
	addr net.Addr
}

// Client enforce http.Client to request to the stub server
// instead of requesting external resources. mapping should be the map from external urls to stub server paths.
func (s *StubServer) Client(mapping map[string]string, c *http.Client) *http.Client {
	stubMapping := make(map[string]string)
	for k, v := range mapping {
		stubMapping[k] = fmt.Sprintf("http://%s%s", s.addr.String(), v)
	}
	return Stub(stubMapping, c)
}

// UseStubServer launches a stub server configured by handler
// and execute test function f.
func UseStubServer(handler http.Handler, f func(*StubServer)) {
	server := &http.Server{}
	listener, err := net.Listen("tcp", "localhost:0")
	xerrors.MustNil(err)
	defer func() {
		listener.Close()
	}()
	server.Handler = handler
	go server.Serve(listener)
	f(&StubServer{
		addr: listener.Addr(),
	})
}
