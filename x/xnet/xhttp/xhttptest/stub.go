package xhttptest

import (
	"fmt"
	"net/http"
	"net/url"
)

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
		return nil, fmt.Errorf("forbitten by xhttptest.Rewriter")
	}
	req.URL = rewriteTo
	return r.base.RoundTrip(req)
}

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
