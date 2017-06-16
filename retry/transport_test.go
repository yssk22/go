package retry

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"io/ioutil"

	"github.com/speedland/go/x/xnet/xhttp/xhttptest"
	"github.com/speedland/go/x/xtesting/assert"
)

func Test_HTTPTransport_SuccessWithRetry(t *testing.T) {
	a := assert.New(t)
	prepareStubServer(
		HTTPAnd(
			HTTPMaxRetry(15),
			HTTPServerErrorChecker(),
		),
		HTTPConstBackoff(100*time.Millisecond),
		func(c *http.Client) {
			resp, err := c.Get("http://example.com/")
			a.Nil(err)
			a.EqInt(200, resp.StatusCode)
			body, _ := ioutil.ReadAll(resp.Body)
			a.EqStr("OK i=10", string(body))
		},
	)
}

func Test_HTTPTransport_FailureWithRtries(t *testing.T) {
	a := assert.New(t)
	prepareStubServer(
		HTTPAnd(
			HTTPMaxRetry(3),
			HTTPServerErrorChecker(),
		),
		HTTPConstBackoff(100*time.Millisecond),
		func(c *http.Client) {
			resp, err := c.Get("http://example.com/")
			a.Nil(err)
			a.EqInt(500, resp.StatusCode)
			body, _ := ioutil.ReadAll(resp.Body)
			a.EqStr("Need Retry i=3", string(body))
		},
	)
}

func prepareStubServer(checker HTTPChecker, backoff HTTPBackoff, f func(*http.Client)) {
	var i = 0
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			i++
			if i < 10 {
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf("Need Retry i=%d", i)))
				return
			}
			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf("OK i=%d", i)))
		}),
		func(s *xhttptest.StubServer) {
			client := s.Client(
				map[string]string{
					"http://example.com/": "/",
				},
				&http.Client{
					Transport: NewHTTPTransport(http.DefaultTransport, checker, backoff),
				},
			)
			f(client)
		},
	)
}
