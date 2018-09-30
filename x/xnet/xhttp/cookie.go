package xhttp

import (
	"fmt"
	"net/http"

	"github.com/yssk22/go/x/xcrypto/xhmac"
)

// SignCookie returns a shallow copy of *http.Cookie with adding hmac signature on it's value.
func SignCookie(c *http.Cookie, hmac *xhmac.Base64) *http.Cookie {
	cc := new(http.Cookie)
	*cc = *c
	cc.Value = hmac.SignString(c.Value)
	return cc
}

// UnsignCookie returns a shallow copy of *http.Cookie with validating hmac signature on it's value.
func UnsignCookie(cc *http.Cookie, hmac *xhmac.Base64) (*http.Cookie, error) {
	rawValue, err := hmac.UnsignString(cc.Value)
	if err != nil {
		return nil, fmt.Errorf("Invalid SignedCookie: %v", err)
	}
	c := new(http.Cookie)
	*c = *cc
	c.Value = rawValue
	return c, nil
}
