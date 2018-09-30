package visualrecognition

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/yssk22/go/x/xarchive/xzip"
	"github.com/yssk22/go/x/xnet/xhttp/xhttptest"
	"github.com/yssk22/go/x/xtesting/assert"

	"context"
)

func TestClient_CreateClassifier(t *testing.T) {
	a := assert.New(t)
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("POST", req.Method)
			a.EqStr("/visual-recognition/api/v3/classifiers", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))

			req.ParseMultipartForm(1 << 16)
			a.EqStr("fruits", req.FormValue("name"))
			for _, cls := range []string{"apple", "banana", "orange"} {
				_, header, err := req.FormFile(fmt.Sprintf("%s_positive_examples", cls))
				a.Nil(err)
				a.EqStr(fmt.Sprintf("%s.zip", cls), header.Filename)
			}
			_, header, err := req.FormFile("negative_examples")
			a.Nil(err)
			a.EqStr("negative.zip", header.Filename)
			http.ServeFile(w, req, "./fixtures/classifier.json")
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			apple, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			banana, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			orange, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			negative, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			defer apple.Close()
			defer banana.Close()
			defer orange.Close()
			defer negative.Close()
			resp, err := c.CreateClassifier(context.Background(),
				"fruits",
				map[string]*xzip.Archiver{
					"apple":  xzip.NewArchiver(apple),
					"banana": xzip.NewArchiver(banana),
					"orange": xzip.NewArchiver(orange),
				},
				xzip.NewArchiver(negative),
			)
			a.Nil(err)
			a.NotNil(resp)
		},
	)
}

func TestClient_GetClassifier(t *testing.T) {
	a := assert.New(t)
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("GET", req.Method)
			a.EqStr("/visual-recognition/api/v3/classifiers/fruits_1050835757", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))
			http.ServeFile(w, req, "./fixtures/classifier.json")
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			resp, err := c.GetClassifier(context.Background(), "fruits_1050835757")
			a.Nil(err)
			a.EqStr("fruits_1050835757", resp.ClassifierID)
			a.EqStr("fruits", resp.Name)
			a.EqStr("abc", resp.Owner)
			a.OK(StatusReady == resp.Status)
			a.EqTime(time.Date(
				2016, 5, 6, 21, 13, 19, 426000000, time.UTC,
			), resp.Created)
			a.EqInt(3, len(resp.Classes))
			a.EqStr("apple", resp.Classes[0].Class)
			a.EqStr("banana", resp.Classes[1].Class)
			a.EqStr("orange", resp.Classes[2].Class)
		},
	)
}

func TestClient_UpdateClassifier(t *testing.T) {
	a := assert.New(t)
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("POST", req.Method)
			a.EqStr("/visual-recognition/api/v3/classifiers/fruits_1050835757", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))

			req.ParseMultipartForm(1 << 16)
			for _, cls := range []string{"apple", "banana", "orange"} {
				_, header, err := req.FormFile(fmt.Sprintf("%s_positive_examples", cls))
				a.Nil(err)
				a.EqStr(fmt.Sprintf("%s.zip", cls), header.Filename)
			}
			_, header, err := req.FormFile("negative_examples")
			a.Nil(err)
			a.EqStr("negative.zip", header.Filename)
			http.ServeFile(w, req, "./fixtures/classifier.json")
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			apple, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			banana, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			orange, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			negative, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			defer apple.Close()
			defer banana.Close()
			defer orange.Close()
			defer negative.Close()
			resp, err := c.UpdateClassifier(context.Background(),
				"fruits_1050835757",
				map[string]*xzip.Archiver{
					"apple":  xzip.NewArchiver(apple),
					"banana": xzip.NewArchiver(banana),
					"orange": xzip.NewArchiver(orange),
				},
				xzip.NewArchiver(negative),
			)
			a.Nil(err)
			a.NotNil(resp)
		},
	)
}

func TestClient_ListClassifiers(t *testing.T) {
	a := assert.New(t)
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("GET", req.Method)
			a.EqStr("/visual-recognition/api/v3/classifiers", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))
			http.ServeFile(w, req, "./fixtures/classifiers.json")
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			resp, err := c.ListClassifiers(context.Background(), false)
			a.Nil(err)
			a.EqInt(2, len(resp.Classifiers))
			a.EqStr("satellite", resp.Classifiers[0].Name)
		},
	)
}

func TestClient_DeleteClassifiers(t *testing.T) {
	a := assert.New(t)
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("DELETE", req.Method)
			a.EqStr("/visual-recognition/api/v3/classifiers/cls1", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))
			w.Write([]byte("{}"))
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			resp, err := c.DeleteClassifier(context.Background(), "cls1")
			a.Nil(err)
			a.OK(resp)
		},
	)
}
