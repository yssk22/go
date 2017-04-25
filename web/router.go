package web

import (
	"net/http"

	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xlog"
)

// Router is a http traffic router
type Router struct {
	middleware *handlerPipeline
	routes     map[string][]*route // (method -> []*route) mapping
	option     *Option
}

// NewRouter returns a new *Router
func NewRouter(option *Option) *Router {
	if option == nil {
		option = DefaultOption
	}
	return &Router{
		middleware: &handlerPipeline{},
		routes:     make(map[string][]*route),
		option:     option,
	}
}

// Use adds middleware handlers to process on every request before all handlers are processed.
func (r *Router) Use(handlers ...Handler) {
	r.middleware.Append(handlers...)
}

// Get adds handlers for "GET {pattern}" requests
func (r *Router) Get(pattern string, handlers ...Handler) {
	r.addRoute("GET", pattern, handlers...)
}

// Post adds handlers for "POST {pattern}" requests
func (r *Router) Post(pattern string, handlers ...Handler) {
	r.addRoute("POST", pattern, handlers...)
}

// Put adds handlers for "PUT {pattern}" requests
func (r *Router) Put(pattern string, handlers ...Handler) {
	r.addRoute("PUT", pattern, handlers...)
}

// Delete adds handlers for "DELETE {pattern}" requests
func (r *Router) Delete(pattern string, handlers ...Handler) {
	r.addRoute("DELETE", pattern, handlers...)
}

func (r *Router) addRoute(method string, pattern string, handlers ...Handler) {
	var rt *route
	if routes, ok := r.routes[method]; !ok {
		rt = newRoute(method, pattern)
		r.routes[method] = []*route{rt}
	} else {
		// if the pattern already exists in current routes, the new handlers should be merged.
		for _, _rt := range routes {
			if _rt.pattern.source == pattern {
				rt = _rt
				break
			}
		}
		if rt == nil {
			rt = newRoute(method, pattern)
			r.routes[method] = append(routes, rt)
		}
	}
	rt.pipeline.Append(handlers...)
}

// Dispatch dispaches *http.Request to the matched handlers and return Response
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	const RequestIDHeader = "X-SPEEDLAND-REQUEST-ID"
	var request = NewRequest(req, r.option)
	var logger = xlog.WithKey("web.router").WithContext(request.Context())
	w.Header().Set(RequestIDHeader, request.ID.String())

	// middleware always executed
	res := r.middleware.Process(
		request,
		NextHandler(func(request *Request) *response.Response {
			// then find a route to dispatch
			path := req.URL.EscapedPath()
			method := req.Method
			route, pathParams := r.findRoute(method, path)
			// bind common fields with request
			if route == nil {
				// Debugging for the route is collectly configured or not.
				logger.Debug(func(p *xlog.Printer) {
					p.Printf("No route is found for \"%s %s\":\n", method, path)
					for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
						for _, r := range r.routes[method] {
							p.Printf("\t%s %s\n", r.method, r.pattern.source)
						}
					}
				})
				return NotFound
			}
			logger.Debug(func(p *xlog.Printer) {
				p.Printf("Routing: %s => %s\n", req.URL.Path, route.pattern.source)
				for _, name := range route.pattern.paramNames {
					p.Printf("\t%s=%s\n", name, pathParams.GetStringOr(name, ""))
				}
			})
			request.Params = pathParams
			return route.pipeline.Process(request.WithValue(requestContextKey, request), nil)
		}),
	)
	if res == nil {
		logger.Debugf("No response is generated.")
		NotFound.Render(request.Context(), w)
		return
	}
	res.Render(request.Context(), w)
}

func (r *Router) findRoute(method string, path string) (*route, *keyvalue.GetProxy) {
	if methodRoutes, ok := r.routes[method]; ok {
		for _, route := range methodRoutes {
			if params, ok := route.pattern.Match(path); ok {
				return route, params
			}
		}
	}
	return nil, nil
}

type route struct {
	method   string
	pattern  *PathPattern
	pipeline *handlerPipeline
}

func newRoute(method, pattern string) *route {
	return &route{
		method:   method,
		pattern:  MustCompilePathPattern(pattern),
		pipeline: &handlerPipeline{},
	}
}
