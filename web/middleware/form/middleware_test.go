package form

import (
	"net/url"
	"testing"

	"github.com/speedland/go/validator"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/httptest"
	"github.com/speedland/go/web/middleware/session"
	"github.com/speedland/go/web/response"
)

var sessionMiddleware = session.NewMiddleware()

func TestMiddleware(t *testing.T) {
	middleware := NewMiddleware(func(v *validator.FormValidator) {
		v.IntField("param").Required().Min(10)
	})
	a := httptest.NewAssert(t)
	recorder := httptest.NewRecorder(prepareRouter(middleware))

	res := recorder.TestPost("/form", nil)
	a.Status(response.HTTPStatusBadRequest, res)
	var e validator.ValidationError
	a.JSON(&e, res)
	a.EqStr("must be required", e.Errors["param"][0].String())

	res = recorder.TestPost("/form", url.Values{
		"param": []string{"0"},
	})
	a.Status(response.HTTPStatusBadRequest, res)
	a.JSON(&e, res)
	a.EqStr("must be more than or equal to 10", e.Errors["param"][0].String())

	res = recorder.TestPost("/form", url.Values{
		"param": []string{"10"},
	})
	a.Status(response.HTTPStatusOK, res)
	a.Body("ok", res)
}

func prepareRouter(middleware *Middleware) *web.Router {
	router := web.NewRouter(nil)
	router.Post("/form",
		middleware,
		web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
			return response.NewText("ok")
		}))
	return router
}
