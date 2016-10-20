package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/speedland/go/web/response"
)

type testHandler (func(req *Request) *response.Text)

func (t testHandler) Process(req *Request) Response {
	return t(req)
}

func ExampleRouter() {
	router := NewRouter()
	router.GET("/path/to/:page.html",
		testHandler(func(req *Request) *response.Text {
			return response.NewText(req.params.GetStringOr("page", ""))
		}),
	)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/path/to/bar.html", nil)
	router.Dispatch(w, req)
	fmt.Printf("Response: %q", w.Body)
	// Output:
	// Response: "bar"
}

func ExampleRouter_multiple_handlers() {
	router := NewRouter()
	router.GET("/path/to/:page.html",
		testHandler(func(req *Request) *response.Text {
			return response.NewText(req.params.GetStringOr("page", ""))
		}),
		testHandler(func(req *Request) *response.Text {
			return response.NewText(req.params.GetStringOr("page", ""))
		}),
	)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/path/to/bar.html", nil)
	router.Dispatch(w, req)
	fmt.Printf("Response: %q", w.Body)
	// Output:
	// Response: "bar"
}
