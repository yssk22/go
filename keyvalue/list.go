package keyvalue

import "fmt"

// List is a list of key-value store.
type List struct {
	list []Getter
	*GetProxy
}

// NewList returns a new *List for `g`
func NewList(g ...Getter) *List {
	list := &List{
		list: g,
	}
	list.GetProxy = NewGetProxy(list) // self
	return list
}

// Get try to get the value from the head of list.
//
//   - If a Getter item returns the value, it is returned.
//   - If a Getter item returns an error other than KeyError, it fails and return that error immediately.
//   - If a Getter item returns KeyError, it tries the next Getter item.
//
func (l *List) Get(key interface{}) (interface{}, error) {
	for _, getter := range l.list {
		v, e := getter.Get(key)
		if e != nil {
			if _, ok := e.(KeyError); !ok {
				return nil, e
			}
		}
		if v != nil {
			return v, nil
		}
	}
	return nil, KeyError(fmt.Sprintf("%s", key))
}
