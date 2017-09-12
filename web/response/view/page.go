package view

import (
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
)

// Page is a struct that configure page endpoints.
type Page interface {
	Render(req *web.Request) *response.Response
}

// PageFunc is to convert a handler function to view.Page
type PageFunc func(*web.Request) Page

// Render implements Page.Render
func (f PageFunc) Render(req *web.Request) *response.Response {
	page := f(req)
	return page.Render(req)
}
