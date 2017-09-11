package view

import (
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
)

// Mount mounts a route of `router` with Page
func Mount(router *web.Router, path string, p Page) {
	router.Get(path, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return p.Render(req)
	}))
}
