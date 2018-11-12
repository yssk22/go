package builtin

import (
	"encoding/json"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/api"
	"github.com/yssk22/go/web/response"
)

func SetupAPI(r web.Router) {
	r.Get("/admin/api/configs/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx := req.Context()
		obj, err := listConfigs(
			ctx,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	r.Get("/admin/api/configs/:key.json", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx := req.Context()
		obj, err := getConfig(
			ctx,
			req.Params.GetStringOr("key", ""),
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	var _updateConfigParameterParser api.ParameterParser
	json.Unmarshal(
		[]byte(`{"specs":{"value":{"type":"string","required":false}},"format":"form"}`),
		&_updateConfigParameterParser,
	)
	r.Put("/admin/api/configs/:key.json", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		var sp updateConfigParams
		if err := _updateConfigParameterParser.Parse(req.Request, &sp); err != nil {
			return err.ToResponse()
		}
		ctx := req.Context()
		obj, err := updateConfig(
			ctx,
			req.Params.GetStringOr("key", ""),
			&sp,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	r.Get("/admin/api/tasks/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx := req.Context()
		obj, err := listAsyncTasks(
			ctx,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	var _createFacebookAuthParameterParser api.ParameterParser
	json.Unmarshal(
		[]byte(`{"specs":{"access_token":{"type":"string","required":false}},"format":"form"}`),
		&_createFacebookAuthParameterParser,
	)
	r.Post("/auth/login/facebook/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		var sp fbLoginParams
		if err := _createFacebookAuthParameterParser.Parse(req.Request, &sp); err != nil {
			return err.ToResponse()
		}
		ctx := req.Context()
		obj, err := createFacebookAuth(
			ctx,
			&sp,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	r.Get("/auth/logout/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx := req.Context()
		obj, err := processLoggedOutForDev(
			ctx,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	r.Post("/auth/logout/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx := req.Context()
		obj, err := processLoggedOut(
			ctx,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
	r.Get("/auth/me.json", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx := req.Context()
		obj, err := getMe(
			ctx,
		)
		if err != nil {
			return api.NewErrorResponse(err)
		}
		return response.NewJSON(obj)
	}))
}
