package example

import (
	"testing"

	"github.com/yssk22/go/web/api/generator/example/types"
	tt "github.com/yssk22/go/types"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/api"
	"github.com/yssk22/go/web/httptest"
	"github.com/yssk22/go/web/response"
)

func Test_getExample(t *testing.T) {
	a := httptest.NewAssert(t)
	router := web.NewRouter(web.DefaultOption)
	SetupAPI(router)

	recorder := httptest.NewRecorder(router)
	resp := recorder.TestGet("/path/to/example/p1/p2/")
	a.Status(response.HTTPStatusOK, resp)
	var data types.ResponseData
	a.JSON(&data, resp)
	a.EqStr("p1", data.StrVal)
}

func Test_getExampleWithExtraParam_Required(t *testing.T) {
	a := httptest.NewAssert(t)
	router := web.NewRouter(web.DefaultOption)
	SetupAPI(router)
	recorder := httptest.NewRecorder(router)

	resp := recorder.TestGet("/path/to/example/p1/2/")
	a.Status(response.HTTPStatusBadRequest, resp)
	var err api.Error
	var fieldErrors api.FieldErrorCollection
	a.JSON(&err, resp)
	a.EqStr("invalid_parameter", err.Code)
	a.EqStr("one or more parameters are invalid", err.Message)
	tt.Typed(err.Extra.(map[string]interface{})["errors"], &fieldErrors)
	a.EqInt(2, len(fieldErrors))
	a.EqInt(1, len(fieldErrors["str_ptr_required"]))
	a.EqInt(1, len(fieldErrors["str_val_required"]))
	a.EqStr("required", fieldErrors["str_ptr_required"][0])
	a.EqStr("required", fieldErrors["str_val_required"][0])
}

func Test_getExampleWithExtraParam(t *testing.T) {
	a := httptest.NewAssert(t)
	router := web.NewRouter(web.DefaultOption)
	SetupAPI(router)

	recorder := httptest.NewRecorder(router)

	// test default
	resp := recorder.TestGet("/path/to/example/p1/2/?str_val_required=rfoo&str_ptr_required=rbar")
	a.Status(response.HTTPStatusOK, resp)
	var data types.RequestParams
	a.JSON(&data, resp)
	a.EqStr("rfoo", data.StrValRequired)
	a.EqStr("foo", data.StrValDefault)
	a.EqStr("", data.StrVal)
	a.EqStr("rbar", *data.StrPtrRequired)
	a.EqStr("bar", *data.StrPtrDefault)
	a.Nil(data.StrPtr)

	// default can be overwritten
	resp = recorder.TestGet("/path/to/example/p1/2/?str_val_required=rfoo&str_ptr_required=rbar&str_val_default=aa&str_ptr_default=bb&str_val=12&str_ptr=34")
	a.Status(response.HTTPStatusOK, resp)
	a.JSON(&data, resp)
	a.EqStr("rfoo", data.StrValRequired)
	a.EqStr("aa", data.StrValDefault)
	a.EqStr("12", data.StrVal)
	a.EqStr("rbar", *data.StrPtrRequired)
	a.EqStr("bb", *data.StrPtrDefault)
	a.EqStr("34", *data.StrPtr)
}

func Test_getExampleWithExtraParamWithValidation(t *testing.T) {
	a := httptest.NewAssert(t)
	router := web.NewRouter(web.DefaultOption)
	SetupAPI(router)

	recorder := httptest.NewRecorder(router)

	// test default
	resp := recorder.TestGet("/path/to/example/p1/2/?str_val_required=rfoo&str_ptr_required=rbar&int_val=-1")
	a.Status(response.HTTPStatusBadRequest, resp)
	var err api.Error
	var fieldErrors api.FieldErrorCollection
	a.JSON(&err, resp)
	a.EqStr("invalid_parameter", err.Code)
	a.EqStr("one or more parameters are invalid", err.Message)
	tt.Typed(err.Extra.(map[string]interface{})["errors"], &fieldErrors)
	a.EqInt(1, len(fieldErrors))
	a.EqInt(1, len(fieldErrors["int_val"]))
	a.EqStr("must be more than or equal to 0", fieldErrors["int_val"][0])
}

// func Test_getExampleWithExtraParam_400(t *testing.T) {
// 	a := httptest.NewAssert(t)
// 	router := web.NewRouter(web.DefaultOption)
// 	SetupAPI(router)

// 	recorder := httptest.NewRecorder(router)
// 	resp := recorder.TestGet("/path/to/example/p1/2/")
// 	a.Status(response.HTTPStatusBadRequest, resp)
// 	var data api.Error
// 	a.JSON(&data, resp)
// 	a.EqStr("invalid_parameter", data.Code)
// }
