package xzip

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/speedland/go/x/xnet/xhttp/xhttptest"
	"github.com/speedland/go/x/xtesting/assert"
)

func TestArchiver(t *testing.T) {
	a := assert.New(t)
	source, err := NewRawSourceFromFile("./fixtures/data.txt")
	a.Nil(err)
	defer source.Close()
	archiver := NewArchiver(source)
	gotBuff, err := ioutil.ReadAll(archiver)
	a.Nil(err)

	expectBuff, _ := ioutil.ReadFile("./fixtures/data.zip")
	a.EqInt(len(expectBuff), len(gotBuff))
	for i := range gotBuff {
		if expectBuff[i] != gotBuff[i] {
			t.Fatalf("Archiver doesn't match expected zip content")
		}
	}
}

func TestNewRawSourceFromFile(t *testing.T) {
	a := assert.New(t)
	source, err := NewRawSourceFromFile("./fixtures/data.txt")
	a.Nil(err)
	defer source.Close()
	a.EqStr("data.txt", source.Name)
	a.EqInt64(19, int64(source.Size))
}

func TestNewRawSourceFromFile_not_exist(t *testing.T) {
	a := assert.New(t)
	_, err := NewRawSourceFromFile("./fixtures/foo.jpg")
	a.NotNil(err)
	a.EqStr(`stat ./fixtures/foo.jpg: no such file or directory`, err.Error())
}

func TestNewRawSourceFromURL(t *testing.T) {
	a := assert.New(t)
	prepareStubServer(func(client *http.Client) {
		source, err := NewRawSourceFromURL("http://example.com/data.txt", client)
		a.Nil(err)
		defer source.Close()
		a.EqStr("data.txt", source.Name)
		a.EqInt64(19, int64(source.Size))
	})
}

func TestNewRawSourceFromURL_non_200(t *testing.T) {
	a := assert.New(t)
	prepareStubServer(func(client *http.Client) {
		_, err := NewRawSourceFromURL("http://example.com/foo", client)
		a.NotNil(err)
		a.EqStr(`non-200 response (404)`, err.Error())
	})
}

func prepareStubServer(f func(*http.Client)) {
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/data.txt" {
				http.ServeFile(w, r, "./fixtures/data.txt")
				return
			}
			if r.URL.Path == "/text.txt" {
				w.Header().Set("content-type", "text/plain; charset=utf-8")
				w.Write([]byte("OK"))
				return
			}
			w.WriteHeader(404)
			w.Write([]byte("not found"))
		}),
		func(s *xhttptest.StubServer) {
			client := s.Client(
				map[string]string{
					"http://example.com/data.txt": "/data.txt",
					"http://example.com/text.txt": "/text.txt",
					"http://example.com/foo":      "/foo",
				},
				&http.Client{},
			)
			f(client)
		},
	)
}
