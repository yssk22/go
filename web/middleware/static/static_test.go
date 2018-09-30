package static

import (
	"testing"

	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/httptest"
)

func TestServeFile(t *testing.T) {
	a := httptest.NewAssert(t)
	router := web.NewRouter(nil)
	router.Get("/", ServeFile("./fixtures/index.html"))
	recorder := httptest.NewRecorder(router)
	res := recorder.TestGet("/")
	a.Body("<html><body>OK</body></html>", res)
}
