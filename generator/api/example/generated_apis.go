package example

import (
	"encoding/json"
	"github.com/yssk22/go/generator/api/example/types"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/api"
	"github.com/yssk22/go/web/response"
)

func SetupAPI(r *web.Router) {
	var _deleteExampleParameterParser api.ParameterParser
	json.Unmarshal(
		[]byte(`{"specs":{"int_val":{"type":"int","required":false},"str_ptr":{"type":"string","required":false},"str_ptr_default":{"type":"string","default":"bar","required":false},"str_ptr_required":{"type":"string","required":true},"str_val":{"type":"string","required":false},"str_val_default":{"type":"string","default":"foo","required":false},"str_val_required":{"type":"string","required":true}},"format":"query"}`),
		&_deleteExampleParameterParser,
	)
	r.Delete("/path/to/example/:param/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		var sp types.RequestParams
		if err := _deleteExampleParameterParser.Parse(req.Request, &sp); err != nil {
			return err.ToResponse()
		}
		ctx := req.Context()
		obj, err := deleteExample(
			ctx,
			req.Params.GetStringOr("param", ""),
			&sp,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	var _createExampleParameterParser api.ParameterParser
	json.Unmarshal(
		[]byte(`{"specs":{"id":{"type":"string","required":false}},"format":"query"}`),
		&_createExampleParameterParser,
	)
	r.Post("/path/to/example/:param/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		var sp Example
		if err := _createExampleParameterParser.Parse(req.Request, &sp); err != nil {
			return err.ToResponse()
		}
		ctx := req.Context()
		obj, err := createExample(
			ctx,
			req.Params.GetStringOr("param", ""),
			&sp,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	var _updateExampleParameterParser api.ParameterParser
	json.Unmarshal(
		[]byte(`{"specs":{"id":{"type":"string","required":false}},"format":"query"}`),
		&_updateExampleParameterParser,
	)
	r.Put("/path/to/example/:param/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		var sp Example
		if err := _updateExampleParameterParser.Parse(req.Request, &sp); err != nil {
			return err.ToResponse()
		}
		ctx := req.Context()
		obj, err := updateExample(
			ctx,
			req.Params.GetStringOr("param", ""),
			&sp,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	var _getExampleWithExtraParamParameterParser api.ParameterParser
	json.Unmarshal(
		[]byte(`{"specs":{"int_val":{"type":"int","required":false},"str_ptr":{"type":"string","required":false},"str_ptr_default":{"type":"string","default":"bar","required":false},"str_ptr_required":{"type":"string","required":true},"str_val":{"type":"string","required":false},"str_val_default":{"type":"string","default":"foo","required":false},"str_val_required":{"type":"string","required":true}},"format":"query"}`),
		&_getExampleWithExtraParamParameterParser,
	)
	r.Get("/path/to/example/:param/2/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		var sp types.RequestParams
		if err := _getExampleWithExtraParamParameterParser.Parse(req.Request, &sp); err != nil {
			return err.ToResponse()
		}
		ctx := req.Context()
		obj, err := getExampleWithExtraParam(
			ctx,
			req.Params.GetStringOr("param", ""),
			&sp,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	r.Get("/path/to/example/:param/:param2/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx := req.Context()
		obj, err := getExample(
			ctx,
			req.Params.GetStringOr("param", ""),
			req.Params.GetStringOr("param2", ""),
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	var _createExample2ParameterParser api.ParameterParser
	json.Unmarshal(
		[]byte(`{"specs":{"inner":{"type":"object","required":false},"int_array":{"type":"array","required":false},"my_enum":{"type":"int","required":false},"str":{"type":"string","required":false}},"format":"json"}`),
		&_createExample2ParameterParser,
	)
	r.Post("/path/to/example2/:param/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		var sp types.ComplexRequestParams
		if err := _createExample2ParameterParser.Parse(req.Request, &sp); err != nil {
			return err.ToResponse()
		}
		ctx := req.Context()
		obj, err := createExample2(
			ctx,
			req.Params.GetStringOr("param", ""),
			&sp,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
}
