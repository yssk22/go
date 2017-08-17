package form

import (
	"net/url"
	"testing"

	"github.com/speedland/go/validator"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/httptest"
	"github.com/speedland/go/web/response"
)

func TestMiddleware_ErrorRequired(t *testing.T) {
	a := httptest.NewAssert(t)
	recorder := httptest.NewRecorder(prepareRouter("param"))
	res := recorder.TestPost("/form", nil)
	a.Status(response.HTTPStatusBadRequest, res)
	var e validator.ValidationError
	a.JSON(&e, res)
	a.EqStr("must be required", e.Errors["param"][0].String())
}

func TestMiddleware_ErrorMin(t *testing.T) {
	a := httptest.NewAssert(t)
	recorder := httptest.NewRecorder(prepareRouter("param"))
	res := recorder.TestPost("/form", url.Values{
		"param": []string{"0"},
	})
	a.Status(response.HTTPStatusBadRequest, res)
	var e validator.ValidationError
	a.JSON(&e, res)
	a.EqStr("must be more than or equal to 10", e.Errors["param"][0].String())
}

func TestMiddleware(t *testing.T) {
	a := httptest.NewAssert(t)
	recorder := httptest.NewRecorder(prepareRouter("param"))
	// should pass validation and response the form value
	res := recorder.TestPost("/form", url.Values{
		"param": []string{"10"},
	})
	a.Status(response.HTTPStatusOK, res)
	a.Body("10", res)
}

func prepareRouter(key string) *web.Router {
	router := web.NewRouter(nil)
	router.Post("/form",
		NewMiddleware(func(v *validator.FormValidator) {
			v.IntField(key).Required().Min(10)
		}),
		web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
			return response.NewText(req.Form.GetStringOr(key, ""))
		}))
	return router
}
