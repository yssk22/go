package keyvalue

// GetProxy provies the proxy functions for Getter to use Get* package functions.
// When you implements Getter interface, you can embed GetProxy to proxy itself to provide
// Get* functions on that struct.
//
//    type MyGetter struct {
//       ...
//       *keyvalue.GetProxy // To
//    }
//
//    func NewMyGetter() *MyGetter {
//       g := &MyGetter{}
//       ...
//       g.GetProxy = keyvalue.GetProxy(g)  // proxy self.
//       return g
//    }
//
type GetProxy struct {
	g Getter
}

// NewGetProxy returns a new GetProxy for Getter
func NewGetProxy(g Getter) *GetProxy {
	return &GetProxy{
		g,
	}
}

// GetOr is shorthand for keyvalue.GetOr.
func (p *GetProxy) GetOr(key string, or interface{}) interface{} {
	return GetOr(p.g, key, or)
}

// GetStringOr is shorthand for keyvalue.GetStringOr.
func (p *GetProxy) GetStringOr(key string, or string) string {
	return GetStringOr(p.g, key, or)
}

// GetIntOr is shorthand for keyvalue.GetIntOr.
func (p *GetProxy) GetIntOr(key string, or int) int {
	return GetIntOr(p.g, key, or)
}
