package example

import (
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
)

func SetupAPI(r *web.Router) {
	r.Get("/path/to/example", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		obj, err := getExample(req)
		if err != nil {
			return response.NewJSON(err)
		}
		return response.NewJSON(obj)
	}))
	r.Get("/path/to/example2", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		obj, err := getExample2(req)
		if err != nil {
			return response.NewJSON(err)
		}
		return response.NewJSON(obj)
	}))
}
