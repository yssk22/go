package web

import (
	"net/url"

	"github.com/speedland/go/keyvalue"
)

func newURLValuesProxy(v url.Values) *keyvalue.GetProxy {
	return keyvalue.NewGetProxy(urlValuesGetProxy(v))
}

type urlValuesGetProxy url.Values

func (p urlValuesGetProxy) Get(key interface{}) (interface{}, error) {
	skey := key.(string)
	if v, ok := p[skey]; ok {
		return v, nil
	}
	return nil, keyvalue.KeyError(skey)
}
