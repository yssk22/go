package web

import (
	"fmt"
	"net/http"

	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/x/xcrypto/xhmac"
	"github.com/speedland/go/x/xnet/xhttp"
	"github.com/speedland/go/x/xtime"
)

func newSignedCookieProxy(cookies []*http.Cookie, hmac *xhmac.Base64) *keyvalue.GetProxy {
	store := make(map[interface{}]*http.Cookie)
	for _, cc := range cookies {
		c, err := xhttp.UnsignCookie(cc, hmac)
		if err == nil {
			store[c.Name] = c
		}
	}
	return keyvalue.NewGetProxy(
		&cookieProxy{
			store: store,
		},
	)
}

type cookieProxy struct {
	store map[interface{}]*http.Cookie
}

func (p *cookieProxy) Get(key interface{}) (interface{}, error) {
	var ok bool
	var v *http.Cookie
	if v, ok = p.store[key]; !ok {
		return nil, keyvalue.KeyError(fmt.Sprintf("%s", key))
	}
	if !v.Expires.IsZero() && v.Expires.After(xtime.Now()) {
		return nil, keyvalue.KeyError(fmt.Sprintf("%s", key))
	}
	return v.Value, nil

}
