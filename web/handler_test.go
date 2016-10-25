package web

import (
	"fmt"
	"net/http"
	"testing"

	"golang.org/x/net/context"

	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xtesting/assert"
)

func Test_handlerPipeline(t *testing.T) {
	a := assert.New(t)
	const contextKey = "test"
	pipeline := &handlerPipeline{}
	pipeline.Append(
		HandlerFunc(func(req *Request, n NextHandler) *response.Response {
			ctx := context.WithValue(req.Context(), contextKey, 1)
			return n(req.WithContext(ctx))
		}),
		HandlerFunc(func(req *Request, n NextHandler) *response.Response {
			i := req.Context().Value(contextKey).(int) + 1
			ctx := context.WithValue(req.Context(), contextKey, i)
			return n(req.WithContext(ctx))
		}),
		HandlerFunc(func(req *Request, n NextHandler) *response.Response {
			i := req.Context().Value(contextKey).(int)
			return response.NewText(fmt.Sprintf("%d", i))
		}),
	)
	r, _ := http.NewRequest("GET", "/", nil)
	res := pipeline.Process(
		NewRequest(r, nil),
		nil,
	)
	a.NotNil(res)
	a.EqStr("2", res.Content())
}

func Test_handlerPipeline_returnNil(t *testing.T) {
	a := assert.New(t)
	pipeline := &handlerPipeline{}
	pipeline.Append(
		HandlerFunc(func(req *Request, n NextHandler) *response.Response {
			return nil
		}),
	)
	r, _ := http.NewRequest("GET", "/", nil)
	res := pipeline.Process(
		NewRequest(r, nil),
		nil,
	)
	a.Nil(res)
}

func Test_handlerPipeline_Multi(t *testing.T) {
	a := assert.New(t)
	const contextKey = "test"
	pipeline1 := &handlerPipeline{}
	pipeline2 := &handlerPipeline{}
	pipeline1.Append(
		HandlerFunc(func(req *Request, n NextHandler) *response.Response {
			ctx := context.WithValue(req.Context(), contextKey, 1)
			return n(req.WithContext(ctx))
		}),
		HandlerFunc(func(req *Request, n NextHandler) *response.Response {
			i := req.Context().Value(contextKey).(int) + 1
			ctx := context.WithValue(req.Context(), contextKey, i)
			return n(req.WithContext(ctx))
		}),
	)
	pipeline2.Append(
		HandlerFunc(func(req *Request, n NextHandler) *response.Response {
			i := req.Context().Value(contextKey).(int)
			return response.NewText(fmt.Sprintf("%d", i))
		}),
	)
	r, _ := http.NewRequest("GET", "/", nil)
	res := pipeline1.Process(
		NewRequest(r, nil),
		NextHandler(func(req *Request) *response.Response {
			return pipeline2.Process(req, nil)
		}),
	)
	a.NotNil(res)
	a.EqStr("2", res.Content())
}
