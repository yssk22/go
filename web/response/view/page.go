package view

import (
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
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
