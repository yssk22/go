package value

import (
	"fmt"

	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/httptest"
	"github.com/yssk22/go/web/response"
)

func ExampleNewRequestValue() {
	val := NewRequestValue(func(req *web.Request) (interface{}, error) {
		return req.URL.Query().Encode(), nil
	})
	r := web.NewRouter(nil)
	r.Get("/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		v, _ := val.Eval(req.Context())
		return response.NewText(req.Context(), v)
	}))
	recorder := httptest.NewRecorder(r)
	fmt.Println(recorder.TestGet("/?foo=bar").Body)
	// Output:
	// foo=bar
}

func ExampleNewQueryIntOr() {
	val := NewQueryIntOr("foo", 10)
	r := web.NewRouter(nil)
	r.Get("/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		v, _ := val.Eval(req.Context())
		return response.NewText(req.Context(), v)
	}))
	recorder := httptest.NewRecorder(r)
	fmt.Println(
		recorder.TestGet("/?foo=5").Body,
		recorder.TestGet("/?foo=bar").Body,
		recorder.TestGet("/").Body,
	)
	// Output:
	// 5 10 10
}

func ExampleNewQueryIntInRange() {
	val := NewQueryIntInRange("foo", 5, 10, 8)
	r := web.NewRouter(nil)
	r.Get("/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		v, _ := val.Eval(req.Context())
		return response.NewText(req.Context(), v)
	}))
	recorder := httptest.NewRecorder(r)
	fmt.Println(
		recorder.TestGet("/?foo=5").Body,
		recorder.TestGet("/?foo=10").Body,
		recorder.TestGet("/?foo=2").Body,
		recorder.TestGet("/?foo=13").Body,
		recorder.TestGet("/?foo=bar").Body,
	)
	// Output:
	// 5 10 8 8 8
}

func ExampleNewQueryIntInList() {
	val := NewQueryIntInList("foo", []int{1}, 8)
	r := web.NewRouter(nil)
	r.Get("/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		v, _ := val.Eval(req.Context())
		return response.NewText(req.Context(), v)
	}))
	recorder := httptest.NewRecorder(r)
	fmt.Println(
		recorder.TestGet("/?foo=1").Body,
		recorder.TestGet("/?foo=2").Body,
		recorder.TestGet("/?foo=bar").Body,
	)
	// Output:
	// 1 8 8
}

func ExampleNewQueryStringOr() {
	val := NewQueryStringOr("foo", "bar")
	r := web.NewRouter(nil)
	r.Get("/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		v, _ := val.Eval(req.Context())
		return response.NewText(req.Context(), v)
	}))
	recorder := httptest.NewRecorder(r)
	fmt.Println(
		recorder.TestGet("/?foo=aaa").Body,
		recorder.TestGet("/").Body,
	)
	// Output:
	// aaa bar
}

func ExampleNewQueryStringInList() {
	val := NewQueryStringInList("foo", []string{"bar"}, "default")
	r := web.NewRouter(nil)
	r.Get("/", web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		v, _ := val.Eval(req.Context())
		return response.NewText(req.Context(), v)
	}))
	recorder := httptest.NewRecorder(r)
	fmt.Println(
		recorder.TestGet("/?foo=aaa").Body,
		recorder.TestGet("/?foo=bar").Body,
	)
	// Output:
	// default bar
}
