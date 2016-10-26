// Package xhttp provies utility functions for net/http package
package xhttp

import (
	"fmt"
	"net/http"
	"strings"
)

const slash byte = '/'

// AbsoluteURL returns an absolute url from path
func AbsoluteURL(req *http.Request, path string) string {
	if strings.HasPrefix(path, "http://") {
		return path
	}
	if strings.HasPrefix(path, "https://") {
		return path
	}
	var scheme, domain string
	if len(path) == 0 {
		path = "/"
	} else {
		if path[0] != slash {
			path = path + "/"
		}
	}

	u := req.URL
	if u.Scheme == "" {
		return path
	}

	scheme = u.Scheme
	domain = u.Host

	return fmt.Sprintf("%s://%s%s", scheme, domain, path)
}
