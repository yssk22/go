package xhttptest

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestStubFile(t *testing.T) {
	a := assert.New(t)
	client := StubFile(map[string]string{
		"http://www.example.com/": "./fixtures/stubfile.txt",
	}, &http.Client{})
	resp, err := client.Get("http://www.example.com/")
	a.Nil(err)
	a.EqStr("text/plain; charset=utf-8", resp.Header.Get("content-type"))
	a.EqInt64(2, resp.ContentLength)
	body, err := ioutil.ReadAll(resp.Body)
	a.Nil(err)
	a.EqStr("OK", string(body))
}

func TestUseStubServer(t *testing.T) {
	a := assert.New(t)
	UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.URL.Path))
		}),
		func(s *StubServer) {
			client := s.Client(map[string]string{
				"http://example.com/foo.txt": "/foo.txt",
			}, &http.Client{})
			resp, err := client.Get("http://example.com/foo.txt")
			a.Nil(err)
			defer resp.Body.Close()
			buff, _ := ioutil.ReadAll(resp.Body)
			a.EqByteString("/foo.txt", buff)

			client = s.Client(nil, &http.Client{})
			resp, err = client.Get("http://example.com/bar.txt")
			a.Nil(err)
			defer resp.Body.Close()
			buff, _ = ioutil.ReadAll(resp.Body)
			a.EqByteString("/bar.txt", buff)
		},
	)
}
