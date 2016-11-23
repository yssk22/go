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

// StubFile is like a Stub but access local file resources
// instead of making actual HTTP requests.
// The following fields are valid and others should not be used in your test.
//
//   - .Status
//   - .StatusCode
//   - .Header.Get("content-type")
//   - .ContentLength
//   - .Body
//
// If you want to test full http request transactions, use StubServer instead.
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
// if mapping is nil, all external requests are mapped to the StubServer s.
func (s *StubServer) Client(mapping map[string]string, c *http.Client) *http.Client {
	if mapping != nil {
		stubMapping := make(map[string]string)
		for k, v := range mapping {
			stubMapping[k] = fmt.Sprintf("http://%s%s", s.addr.String(), v)
		}
		return stubByMap(stubMapping, c)
	}
	return stubByServer(s, c)
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

// stubByMap enforce stub accesses for http.Client using URL mapping.
func stubByMap(mapping map[string]string, c *http.Client) *http.Client {
	rewriteMapping := make(map[string]*url.URL)
	for k, v := range mapping {
		var err error
		rewriteMapping[k], err = url.Parse(v)
		if err != nil {
			panic(fmt.Errorf("invalid stub url mapping: %s -> %s", k, v))
		}
	}
	c.Transport = &mappingRewriter{
		mapping: rewriteMapping,
		base:    c.Transport,
	}
	return c
}

// mappingRewriter is to rewrite a request not to access a remote resource
type mappingRewriter struct {
	mapping map[string]*url.URL
	server  *StubServer
	base    http.RoundTripper
}

// RoundTrip implements http.RoundTripper#RoundTrip
func (r *mappingRewriter) RoundTrip(req *http.Request) (*http.Response, error) {
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

// stubByServer enforce stub accesses for http.Client to the StubServer
func stubByServer(s *StubServer, c *http.Client) *http.Client {
	c.Transport = &serverRewriter{
		server: s,
		base:   c.Transport,
	}
	return c
}

// serverRewriter is to rewrite a request not to access a remote resource
type serverRewriter struct {
	server *StubServer
	base   http.RoundTripper
}

// RoundTrip implements http.RoundTripper#RoundTrip
func (r *serverRewriter) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = r.server.addr.String()
	if r.base != nil {
		return r.base.RoundTrip(req)
	}
	return http.DefaultTransport.RoundTrip(req)
}
