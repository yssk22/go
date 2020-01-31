package form

import (
	"github.com/yssk22/go/validator"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
)

type Middleware struct {
	validator *validator.FormValidator
}

func NewMiddleware(f func(*validator.FormValidator)) *Middleware {
	v := validator.NewFormValidator()
	if f != nil {
		f(v)
	}
	return &Middleware{
		validator: v,
	}
}

func (m *Middleware) Process(req *web.Request, next web.NextHandler) *response.Response {
	req.ParseForm()
	err := m.validator.Eval(req.PostForm)
	if err != nil {
		return response.NewJSONWithStatus(
			req.Context(), err, response.HTTPStatusBadRequest,
		)
	}
	return next(req)
}
