package keyvalue

import (
	"net/url"
)

// NewQueryProxy returns *GetProxy for url.Values
func NewQueryProxy(query url.Values) *GetProxy {
	return GetterStringKeyFunc(func(key string) (interface{}, error) {
		if val := query[key]; val != nil {
			return val[0], nil
		}
		return nil, KeyError(key)
	}).Proxy()
}
