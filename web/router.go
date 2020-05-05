package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/web/response"
	"github.com/yssk22/go/x/xlog"
)

// Router is an interface to set up http router
type Router interface {
	Use(...Handler)
	All(string, ...Handler)
	Get(string, ...Handler)
	Post(string, ...Handler)
	Put(string, ...Handler)
	Delete(string, ...Handler)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// defaultRouter is a http traffic router default implementation
type defaultRouter struct {
	middleware *handlerPipeline
	routes     map[string][]*route // (method -> []*route) mapping
	option     *Option
}

// NewRouter returns a new *Router
func NewRouter(option *Option) Router {
	if option == nil {
		option = DefaultOption
	}
	r := &defaultRouter{
		middleware: &handlerPipeline{},
		routes:     make(map[string][]*route),
		option:     option,
	}

	r.Get("/__debug__/routes", HandlerFunc(func(req *Request, next NextHandler) *response.Response {
		var buff bytes.Buffer
		r.printRoutes("GET", &buff)
		r.printRoutes("POST", &buff)
		r.printRoutes("PUT", &buff)
		r.printRoutes("DELETE", &buff)
		return response.NewText(req.Context(), buff.String())
	}))
	return r
}

func (r *defaultRouter) printRoutes(method string, dst io.Writer) {
	if routes, ok := r.routes[method]; ok {
		for _, r := range routes {
			if _, err := dst.Write(
				[]byte(fmt.Sprintf("%s %s\n", r.method, r.pattern.source)),
			); err != nil {
				panic(err)
			}
		}
	}
}

// Use adds middleware handlers to process on every request before all handlers are processed.
func (r *defaultRouter) Use(handlers ...Handler) {
	r.middleware.Append(handlers...)
}

// All adds handlers for "GET|PUT|POST|DELETE {pattern}" requests
func (r *defaultRouter) All(pattern string, handlers ...Handler) {
	r.addRoute("GET", pattern, handlers...)
	r.addRoute("PUT", pattern, handlers...)
	r.addRoute("POST", pattern, handlers...)
	r.addRoute("DELETE", pattern, handlers...)
}

// Get adds handlers for "GET {pattern}" requests
func (r *defaultRouter) Get(pattern string, handlers ...Handler) {
	r.addRoute("GET", pattern, handlers...)
}

// Post adds handlers for "POST {pattern}" requests
func (r *defaultRouter) Post(pattern string, handlers ...Handler) {
	r.addRoute("POST", pattern, handlers...)
}

// Put adds handlers for "PUT {pattern}" requests
func (r *defaultRouter) Put(pattern string, handlers ...Handler) {
	r.addRoute("PUT", pattern, handlers...)
}

// Delete adds handlers for "DELETE {pattern}" requests
func (r *defaultRouter) Delete(pattern string, handlers ...Handler) {
	r.addRoute("DELETE", pattern, handlers...)
}

func (r *defaultRouter) addRoute(method string, pattern string, handlers ...Handler) {
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
func (r *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	const RequestIDHeader = "X-SPEEDLAND-REQUEST-ID"
	var request = NewRequest(req, r.option)
	var logger = xlog.WithKey("web.router").WithContext(request.Context())
	w.Header().Set(RequestIDHeader, request.ID.String())

	// middleware always executed
	r.middleware.Process(
		request,
		NextHandler(func(request *Request) *response.Response {
			// then find a route to dispatch
			path := req.URL.EscapedPath()
			method := req.Method
			matched := r.findRoute(method, path)
			// bind common fields with request
			if matched == nil {
				// Debugging for the route is collectly configured or not.
				logger.Debug(func(p *xlog.Printer) {
					p.Printf("No route is found for \"%s %s\":\n", method, path)
					for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
						for _, r := range r.routes[method] {
							p.Printf("\t%s %s\n", r.method, r.pattern.source)
						}
					}
				})
				r.renderResponse(request.Context(), w, response.NewTextWithStatus(request.Context(), "not found", response.HTTPStatusNotFound))
				return nil
			}
			if len(matched) == 1 {
				route := matched[0].Route
				pathParams := matched[0].Params
				logger.Debug(func(p *xlog.Printer) {
					p.Printf("Routing matched: %s => %s", req.URL.Path, route.pattern.source)
					for _, name := range route.pattern.paramNames {
						p.Printf("\n\t%s=%s", name, pathParams.GetStringOr(name, ""))
					}
				})
				request.Params = pathParams
				res := route.pipeline.Process(request.WithValue(requestContextKey, request), nil)
				r.renderResponse(req.Context(), w, res)
				return nil
			}
			// dynamic creation of pipeline
			var handlers []Handler
			for _, m := range matched {
				route := m.Route
				pathParams := m.Params
				logger.Debug(func(p *xlog.Printer) {
					p.Printf("Routing: %s => %s\n", req.URL.Path, route.pattern.source)
					for _, name := range route.pattern.paramNames {
						p.Printf("\t%s=%s\n", name, pathParams.GetStringOr(name, ""))
					}
				})
				handlers = append(handlers, HandlerFunc(func(req *Request, next NextHandler) *response.Response {
					req.Params = pathParams
					return next(req)
				}))
				handlers = append(handlers, m.Route.pipeline.Handlers...)
			}
			var dynamic = &handlerPipeline{}
			dynamic.Append(handlers...)
			res := dynamic.Process(request.WithValue(requestContextKey, request), nil)
			r.renderResponse(req.Context(), w, res)
			return nil
		}),
	)
}

func (r *defaultRouter) renderResponse(ctx context.Context, w http.ResponseWriter, res *response.Response) {
	if res == nil {
		response.NewTextWithStatus(ctx, "not found", response.HTTPStatusNotFound).Render(w)
		return
	}
	res.Render(w)
}

type matchedRoute struct {
	Route  *route
	Params *keyvalue.GetProxy
}

func (r *defaultRouter) findRoute(method string, path string) []*matchedRoute {
	if methodRoutes, ok := r.routes[method]; ok {
		var matched []*matchedRoute
		var found = false
		for _, route := range methodRoutes {
			if params, ok := route.pattern.Match(path); ok {
				matched = append(matched, &matchedRoute{
					route,
					params,
				})
				found = true
			}
		}
		if found {
			return matched
		}
	}
	return nil
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
