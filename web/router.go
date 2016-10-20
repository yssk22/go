package web

import "net/http"

// Router is a http traffic router
type Router struct {
	before []Handler
	after  []Handler
	routes map[string][]*route // method -> []Handler
}

// NewRouter returns a new *Router
func NewRouter() *Router {
	return &Router{
		before: make([]Handler, 0),
		after:  make([]Handler, 0),
		routes: make(map[string][]*route),
	}
}

// GET adds handlers for "GET {pattern}" requests
func (r *Router) GET(pattern string, handlers ...Handler) {
	const method = "GET"
	r.addRoute(method, pattern, handlers...)
}

func (r *Router) addRoute(method string, pattern string, handlers ...Handler) {
	rt := &route{
		method:   method,
		pattern:  MustCompilePathPattern(pattern),
		handlers: handlers,
	}
	if _, ok := r.routes[method]; !ok {
		r.routes[method] = make([]*route, 0)
	}
	r.routes[method] = append(r.routes[method], rt)
}

// Dispatch dispaches *http.Request to the matched handlers and return Response
func (r *Router) Dispatch(w http.ResponseWriter, req *http.Request) {
	var request = &Request{
		Request: req,
	}
	// after handlers should always be processed.
	defer func() {
		for _, h := range r.after {
			h.Process(request)
		}
	}()

	// before handlers
	for _, h := range r.before {
		if res := h.Process(request); res != nil {
			if res.Render(w) {
				return
			}
		}
	}

	// main request handlers
	path := req.URL.EscapedPath()
	method := req.Method

	request.query = newURLValuesProxy(req.URL.Query())
	if methodRoutes, ok := r.routes[method]; ok {
		for _, route := range methodRoutes {
			if params, ok := route.pattern.Match(path); ok {
				request.params = params
				for _, h := range route.handlers {
					if res := h.Process(request); res != nil {
						if res.Render(w) {
							return
						}
					}
				}
			}
		}
	}
}

type route struct {
	method   string
	pattern  *PathPattern
	handlers []Handler
}
