// Package value provides lazy.Value for web context.
package value

import (
	"fmt"

	"context"

	"github.com/yssk22/go/lazy"
	"github.com/yssk22/go/web"
)

type lazyRequestValue struct {
	eval func(req *web.Request) (interface{}, error)
}

func (l *lazyRequestValue) Eval(ctx context.Context) (interface{}, error) {
	if req := web.FromContext(ctx); req != nil {
		return l.eval(req)
	}
	return nil, fmt.Errorf("not a request context")
}

// NewRequestValue returns a lazy.Value in a request context
func NewRequestValue(f func(*web.Request) (interface{}, error)) lazy.Value {
	return &lazyRequestValue{
		eval: f,
	}
}

// NewQueryIntOr returns a new *lazy.Value to get the query value by `key` or `or` value if the value is empty or invalid.
func NewQueryIntOr(key string, or int) lazy.Value {
	return NewRequestValue(func(req *web.Request) (interface{}, error) {
		return req.Query.GetIntOr(key, or), nil
	})
}

// NewQueryIntInRange is like NewQueryIntOr but it also limits the min and max for the value
func NewQueryIntInRange(key string, min int, max int, or int) lazy.Value {
	return NewRequestValue(func(req *web.Request) (interface{}, error) {
		if v := req.Query.GetIntOr(key, or); v >= min && v <= max {
			return v, nil
		}
		return or, nil
	})
}

// NewQueryIntInList is like NewQueryIntOr but it also limits the possible values in the `list``
func NewQueryIntInList(key string, values []int, or int) lazy.Value {
	return NewRequestValue(func(req *web.Request) (interface{}, error) {
		v := req.Query.GetIntOr(key, or)
		for _, vv := range values {
			if v == vv {
				return v, nil
			}
		}
		return or, nil
	})
}

// NewQueryStringOr returns a new *lazy.Value to get the query value by `key` or `or` value if the value is empty.
func NewQueryStringOr(key string, or string) lazy.Value {
	return NewRequestValue(func(req *web.Request) (interface{}, error) {
		return req.Query.GetStringOr(key, or), nil
	})
}

// NewQueryStringInList is like NewQueryStringOr but it also limits the possible values in the `list``
func NewQueryStringInList(key string, values []string, or string) lazy.Value {
	return NewRequestValue(func(req *web.Request) (interface{}, error) {
		v := req.Query.GetStringOr(key, or)
		for _, vv := range values {
			if v == vv {
				return v, nil
			}
		}
		return or, nil
	})
}
