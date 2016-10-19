package keyvalue

// List is a list of key-value store.
type List struct {
	list []Getter
}

// NewList returns a new *List for `g`
func NewList(g ...Getter) *List {
	list := &List{
		list: g,
	}
	return list
}

// Get try to get the value from the head of list.
//
//   - If a Getter item returns the value, it is returned.
//   - If a Getter item returns an error other than KeyError, it fails and return that error immediately.
//   - If a Getter item returns KeyError, it tries the next Getter item.
//
func (l *List) Get(key string) (interface{}, error) {
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
	return nil, KeyError(key)
}

// GetOr returns a value from the default config list or a default value if not found.
func (l *List) GetOr(key string, or interface{}) interface{} {
	return GetOr(l, key, or)
}

// GetStringOr is a string version of GetOr
func (l *List) GetStringOr(key string, or string) string {
	return GetStringOr(l, key, or)
}

// GetIntOr is a int version of GetOr
func (l *List) GetIntOr(key string, or int) int {
	return GetIntOr(l, key, or)
}
