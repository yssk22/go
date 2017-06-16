package view

import (
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
)

// Page is a struct that configure page endpoints.
type Page interface {
	Render(req *web.Request) *response.Response
}
