package visualrecognition

import (
	"encoding/json"
	"net/http"
	"testing"

	"context"

	"github.com/yssk22/go/x/xarchive/xzip"
	"github.com/yssk22/go/x/xnet/xhttp/xhttptest"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestClient_ClassifyURL(t *testing.T) {
	a := assert.New(t)
	targetURL := "https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/prez.jpg"
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("GET", req.Method)
			a.EqStr("/visual-recognition/api/v3/classify", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))
			a.EqStr(targetURL, req.URL.Query().Get("url"))
			http.ServeFile(w, req, "./fixtures/classify.json")
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			resp, err := c.ClassifyURL(context.Background(), targetURL, nil)
			a.Nil(err)
			assertClassifyResponse(a, resp)
		},
	)
}

func TestClient_ClassifyURL_with_params(t *testing.T) {
	a := assert.New(t)
	targetURL := "https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/prez.jpg"
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("GET", req.Method)
			a.EqStr("/visual-recognition/api/v3/classify", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))
			a.EqStr("cid1,cid2", req.URL.Query().Get("classifier_ids"))
			a.EqStr("0.8", req.URL.Query().Get("threshold"))
			a.EqStr(targetURL, req.URL.Query().Get("url"))
			http.ServeFile(w, req, "./fixtures/classify.json")
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			params := &ClassifyParams{
				ClassifierIDs: []string{"cid1", "cid2"},
				Threshold:     0.8,
			}
			resp, err := c.ClassifyURL(context.Background(), targetURL, params)
			a.Nil(err)
			assertClassifyResponse(a, resp)
		},
	)
}

func TestClient_ClassifyImages(t *testing.T) {
	a := assert.New(t)
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("POST", req.Method)
			a.EqStr("/visual-recognition/api/v3/classify", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))

			req.ParseMultipartForm(1 << 16)
			_, header, err := req.FormFile("images_file")
			a.Nil(err)
			a.EqStr("image.zip", header.Filename)
			http.ServeFile(w, req, "./fixtures/classify.json")
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			rs, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			defer rs.Close()
			resp, err := c.ClassifyImages(context.Background(), xzip.NewArchiver(rs), nil)
			a.Nil(err)
			assertClassifyResponse(a, resp)
		},
	)
}

func TestClient_ClassifyImages_with_params(t *testing.T) {
	a := assert.New(t)
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("POST", req.Method)
			a.EqStr("/visual-recognition/api/v3/classify", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))

			req.ParseMultipartForm(1 << 16)
			_, header, err := req.FormFile("images_file")
			a.Nil(err)
			a.EqStr("image.zip", header.Filename)

			file, header, err := req.FormFile("parameters")
			a.Nil(err)
			a.EqStr("parameters.json", header.Filename)
			var params ClassifyParams
			a.Nil(json.NewDecoder(file).Decode(&params))
			a.EqInt(2, len(params.ClassifierIDs))
			a.EqStr("cid1", params.ClassifierIDs[0])
			a.EqStr("cid2", params.ClassifierIDs[1])
			a.EqFloat64(0.8, params.Threshold)
			http.ServeFile(w, req, "./fixtures/classify.json")
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			params := &ClassifyParams{
				ClassifierIDs: []string{"cid1", "cid2"},
				Threshold:     0.8,
			}
			rs, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			defer rs.Close()
			resp, err := c.ClassifyImages(context.Background(), xzip.NewArchiver(rs), params)
			a.Nil(err)
			assertClassifyResponse(a, resp)
		},
	)
}

func assertClassifyResponse(a *assert.Assert, resp *ClassifyResponse) {
	a.EqInt(1, len(resp.Images))
	a.EqInt(2, len(resp.Images[0].Classifiers))
	classifier := resp.Images[0].Classifiers[0]
	a.EqInt(6, len(classifier.Classes))
	a.EqStr("apple", classifier.Classes[0].Class)
	a.EqFloat64(0.645656, classifier.Classes[0].Score)
	a.EqStr("default", classifier.ClassifierID)
	a.EqStr("default", classifier.Name)
}

// func TestClient_ClassifyByURL(t *testing.T) {
// 	a := assert.New(t)
// 	stubClient := xhttptest.StubFile(
// 		map[string]string{
// 			"https://gateway-a.watsonplatform.net/visual-recognition/api/v3/classify?api_key=abc&url=https%3A%2F%2Fgithub.com%2Fwatson-developer-cloud%2Fdoc-tutorial-downloads%2Fraw%2Fmaster%2Fvisual-recognition%2Fprez.jpg&version=2016-05-20": "./fixtures/classify.json",
// 		},
// 		&http.Client{},
// 	)
// 	c := NewClient("abc", stubClient)
// 	resp, err := c.ClassifyByURL(context.Background(), "https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/prez.jpg", nil)
// 	a.Nil(err)
// }
