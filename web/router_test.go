package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/speedland/go/web/response"
)

func ExampleRouter() {
	router := NewRouter()
	router.Get("/path/to/:page.html",
		HandlerFunc(func(req *Request, _ NextHandler) Response {
			return response.NewText(req.Params.GetStringOr("page", ""))
		}),
	)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/path/to/bar.html", nil)
	router.Dispatch(w, req)
	fmt.Printf("Response: %q", w.Body)
	// Output:
	// Response: "bar"
}

func ExampleRouter_multipleHandlerPipeline() {
	router := NewRouter()
	router.Get("/path/to/:page.html",
		HandlerFunc(func(req *Request, next NextHandler) Response {
			if req.Params.GetStringOr("page", "") == "first" {
				return response.NewText("First Handler")
			}
			return next(req)
		}),
		HandlerFunc(func(req *Request, next NextHandler) Response {
			// This handler is reached only when the first handler returns nil
			if req.Params.GetStringOr("page", "") == "second" {
				return response.NewText("Second Handler")
			}
			return nil
		}),
	)
	for _, s := range []string{"first", "second"} {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/path/to/%s.html", s), nil)
		router.Dispatch(w, req)
		fmt.Printf("Response: %q\n", w.Body)
	}
	// Output:
	// Response: "First Handler"
	// Response: "Second Handler"
}

func ExampleRouter_Use() {
	router := NewRouter()
	router.Use(HandlerFunc(func(req *Request, next NextHandler) Response {
		return next(req.WithValue(
			"my-middleware-key",
			"my-middleware-value",
		))
	}))
	router.Get("/a.html",
		HandlerFunc(func(req *Request, next NextHandler) Response {
			v, _ := req.Get("my-middleware-key")
			return response.NewText(v.(string))
		}),
	)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/a.html", nil)
	router.Dispatch(w, req)
	fmt.Printf("Response: %q\n", w.Body)

	// Output:
	// Response: "my-middleware-value"
}
