package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yssk22/go/web/response"
	"github.com/yssk22/go/x/xtesting/assert"
)

func ExampleRouter() {
	router := NewRouter(nil)
	router.Get("/path/to/:page.html",
		HandlerFunc(func(req *Request, _ NextHandler) *response.Response {
			return response.NewText(req.Context(), req.Params.GetStringOr("page", ""))
		}),
	)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/path/to/bar.html", nil)
	router.ServeHTTP(w, req)
	fmt.Printf("*response.Response: %q", w.Body)
	// Output:
	// *response.Response: "bar"
}

func TestRouter_multipleRoutes(t *testing.T) {
	a := assert.New(t)
	router := NewRouter(nil)
	router.Get("/a.html",
		HandlerFunc(func(req *Request, _ NextHandler) *response.Response {
			return response.NewText(req.Context(), "a.html")
		}),
	)
	router.Get("/b.html",
		HandlerFunc(func(req *Request, _ NextHandler) *response.Response {
			return response.NewText(req.Context(), "b.html")
		}),
	)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/a.html", nil)
	router.ServeHTTP(w, req)
	a.EqStr("a.html", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/b.html", nil)
	router.ServeHTTP(w, req)
	a.EqStr("b.html", w.Body.String())
}

func ExampleRouter_multipleHandlerPipeline() {
	router := NewRouter(nil)
	router.Get("/path/to/:page.html",
		HandlerFunc(func(req *Request, next NextHandler) *response.Response {
			if req.Params.GetStringOr("page", "") == "first" {
				return response.NewText(req.Context(), "First Handler")
			}
			return next(req)
		}),
		HandlerFunc(func(req *Request, next NextHandler) *response.Response {
			// This handler is reached only when the first handler returns nil
			if req.Params.GetStringOr("page", "") == "second" {
				return response.NewText(req.Context(), "Second Handler")
			}
			return nil
		}),
	)
	for _, s := range []string{"first", "second"} {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/path/to/%s.html", s), nil)
		router.ServeHTTP(w, req)
		fmt.Printf("*response.Response: %q\n", w.Body)
	}
	// Output:
	// *response.Response: "First Handler"
	// *response.Response: "Second Handler"
}

func ExampleRouter_Use() {
	var state []string
	router := NewRouter(nil)
	router.Use(HandlerFunc(func(req *Request, next NextHandler) *response.Response {
		state = append(state, "middleware-handler")
		resp := next(req.WithValue(
			"my-middleware-key",
			"my-middleware-value",
		))
		state = append(state, "middleware-finalize")
		return resp
	}))
	router.Get("/a.html",
		HandlerFunc(func(req *Request, next NextHandler) *response.Response {
			state = append(state, "request-handler")
			v, _ := req.Get("my-middleware-key")
			return response.NewText(req.Context(), v.(string))
		}),
	)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/a.html", nil)
	router.ServeHTTP(w, req)
	fmt.Printf("*response.Response: %q\n", w.Body)
	fmt.Printf("state: %s", strings.Join(state, ">"))
	// Output:
	// *response.Response: "my-middleware-value"
	// state: middleware-handler>request-handler>middleware-finalize
}

func ExampleRouter_multipleRoute() {
	router := NewRouter(nil)
	router.Get("/:key.html", HandlerFunc(func(req *Request, next NextHandler) *response.Response {
		return next(req.WithValue(
			"my-middleware-key",
			req.Params.GetStringOr("key", "default"),
		))
	}))
	router.Get("/a.html",
		HandlerFunc(func(req *Request, next NextHandler) *response.Response {
			v, _ := req.Get("my-middleware-key")
			return response.NewText(req.Context(), fmt.Sprintf("a-%s", v))
		}),
	)
	router.Get("/b.html",
		HandlerFunc(func(req *Request, next NextHandler) *response.Response {
			v, _ := req.Get("my-middleware-key")
			return response.NewText(req.Context(), fmt.Sprintf("b-%s", v))
		}),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/a.html", nil)
	router.ServeHTTP(w, req)
	fmt.Printf("*response.Response: %q\n", w.Body)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/b.html", nil)
	router.ServeHTTP(w, req)
	fmt.Printf("*response.Response: %q\n", w.Body)

	// not found route even /:key.html handles some
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/c.html", nil)
	router.ServeHTTP(w, req)
	fmt.Printf("*response.Response: %q\n", w.Body)

	// Output:
	// *response.Response: "a-a"
	// *response.Response: "b-b"
	// *response.Response: "not found"
}

func TestRouter_middlewareBeforeAfter(t *testing.T) {
	a := assert.New(t)
	var state []string
	router := NewRouter(nil)
	router.Use(HandlerFunc(func(req *Request, next NextHandler) *response.Response {
		state = append(state, "middleware-before")
		resp := next(req)
		state = append(state, "middleware-after")
		return resp
	}))

	router.Get("/", HandlerFunc(func(req *Request, next NextHandler) *response.Response {
		return response.NewHandler(
			req.Context(),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				state = append(state, "app")
			}),
			req.Request,
		)
	}))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	a.EqInt(3, len(state))
	a.EqStr("middleware-before", state[0])
	a.EqStr("app", state[1])
	a.EqStr("middleware-after", state[2])
}
