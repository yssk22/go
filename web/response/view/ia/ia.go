package ia

import (
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
)

type Func func(*web.Request) (*PageVars, error)

// Render implements view.Page#Render
func (f Func) Render(req *web.Request) *response.Response {
	v, err := f(req)
	if err != nil {
		return response.NewTextWithStatus(req.Context(), "internal server error", response.HTTPStatusInternalServerError)
	}
	if v == nil {
		return response.NewTextWithStatus(req.Context(), "no content", response.HTTPStatusNoContent)
	}
	return response.NewHTMLWithStatus(
		req.Context(),
		iaMarkupTemplate,
		v,
		response.HTTPStatusOK,
	)
}
