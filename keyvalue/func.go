package keyvalue

// GetterFunc is to create a proxy wrapper for Getter
type GetterFunc func(interface{}) (interface{}, error)

// Get implements Getter#Get
func (f GetterFunc) Get(key interface{}) (interface{}, error) {
	return f(key)
}

// Proxy returns GetProxy for this func.
func (f GetterFunc) Proxy() *GetProxy {
	return NewGetProxy(f)
}

// GetterStringKeyFunc is to create a proxy wrapper for Getter
type GetterStringKeyFunc func(string) (interface{}, error)

// Get implements Getter#Get
func (f GetterStringKeyFunc) Get(key interface{}) (interface{}, error) {
	return f(key.(string))
}

// Proxy returns GetProxy for this func.
func (f GetterStringKeyFunc) Proxy() *GetProxy {
	return NewGetProxy(f)
}
