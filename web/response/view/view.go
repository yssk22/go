package view

import (
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
)

// Mount mounts a route of `router` with Page
func Mount(router *web.Router, path string, p Page) {
	router.Get(path, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return p.Render(req)
	}))
}
