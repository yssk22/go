package web

import (
	"net/url"

	"github.com/speedland/go/keyvalue"
)

func newURLValuesProxy(v url.Values) *keyvalue.GetProxy {
	return keyvalue.NewGetProxy(urlValuesGetProxy(v))
}

type urlValuesGetProxy url.Values

func (p urlValuesGetProxy) Get(key string) (interface{}, error) {
	if v, ok := p[key]; ok {
		return v, nil
	}
	return nil, keyvalue.KeyError(key)
}
