package cors

import (
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
)

type middleware struct {
	m         map[string]bool
	all       bool
	preflight web.Handler
}

var defaultPreflight = web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
	resp := response.NewText("OK")
	return resp
})

// NewMiddleware returns a web.Handler interface for CORS support
func NewMiddleware(origins ...string) web.Handler {
	return NewMiddlewareWithPreflight(defaultPreflight, origins...)
}

// NewMiddlewareWithPreflight returns a web.Handler interface for CORS support
// by specifying the preflight request handler
func NewMiddlewareWithPreflight(preflight web.Handler, origins ...string) web.Handler {
	var m = &middleware{
		m:         make(map[string]bool),
		all:       false,
		preflight: preflight,
	}
	for _, o := range origins {
		m.m[o] = true
	}
	if _, ok := m.m["*"]; ok {
		m.all = true
	}
	return m
}

func (m *middleware) Process(req *web.Request, next web.NextHandler) *response.Response {
	var resp *response.Response
	if req.Method == "OPTIONS" {
		resp = m.preflight.Process(req, next)
		resp.Header.Add("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,HEAD,OPTIONS")
		resp.Header.Add("Access-Control-Allow-Headers", "Origin,Authorization,Accept,Content-Type")
	} else {
		resp = next(req)
	}
	if resp != nil {
		if m.all {
			resp.Header.Add("Access-Control-Allow-Origin", "*")
		} else {
			origin := req.Request.Header.Get("Origin")
			if origin != "" {
				if _, ok := m.m[origin]; ok {
					resp.Header.Add("Access-Control-Allow-Origin", origin)
				}
			}
		}
	}
	return resp
}
