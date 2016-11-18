package static

import (
	"testing"

	"github.com/speedland/go/web"
	"github.com/speedland/go/web/httptest"
)

func TestServeFile(t *testing.T) {
	a := httptest.NewAssert(t)
	router := web.NewRouter(nil)
	router.Get("/", ServeFile("./fixtures/index.html"))
	recorder := httptest.NewRecorder(router)
	res := recorder.TestGet("/")
	a.Body("<html><body>OK</body></html>", res)
}
