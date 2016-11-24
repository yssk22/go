package visualrecognition

import (
	"net/http"
	"testing"

	"github.com/speedland/go/x/xarchive/xzip"
	"github.com/speedland/go/x/xnet/xhttp/xhttptest"
	"github.com/speedland/go/x/xtesting/assert"
	"golang.org/x/net/context"
)

func TestClient_DetectFacesOnURL(t *testing.T) {
	a := assert.New(t)
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("GET", req.Method)
			a.EqStr("/visual-recognition/api/v3/detect_faces", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))
			a.EqStr("https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/prez.jpg", req.URL.Query().Get("url"))
			http.ServeFile(w, req, "./fixtures/detect-face.json")
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			resp, err := c.DetectFacesOnURL(context.Background(), "https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/prez.jpg")
			a.Nil(err)
			assertFaceDetectResponse(a, resp)
		},
	)
}

func TestClient_DetectFacesOnImages(t *testing.T) {
	a := assert.New(t)
	xhttptest.UseStubServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			a.EqStr("POST", req.Method)
			a.EqStr("/visual-recognition/api/v3/detect_faces", req.URL.Path)
			a.EqStr("abc", req.URL.Query().Get("api_key"))
			a.EqStr("2016-05-20", req.URL.Query().Get("version"))

			req.ParseMultipartForm(1 << 16)
			_, header, err := req.FormFile("images_file")
			a.Nil(err)
			a.EqStr("image.zip", header.Filename)
			http.ServeFile(w, req, "./fixtures/detect-face.json")
		}),
		func(s *xhttptest.StubServer) {
			c := NewClient("abc", s.Client(nil, &http.Client{}))
			rs, _ := xzip.NewRawSourceFromFile("./fixtures/prez.jpg")
			defer rs.Close()
			resp, err := c.DetectFacesOnImages(context.Background(), xzip.NewArchiver(rs))
			a.Nil(err)
			assertFaceDetectResponse(a, resp)
		},
	)
}

func assertFaceDetectResponse(a *assert.Assert, resp *FaceDetectResponse) {
	a.EqInt(1, len(resp.Images))
	a.EqInt(1, len(resp.Images[0].Faces))
	face := resp.Images[0].Faces[0]
	a.NotNil(face.Age)
	a.EqInt(54, face.Age.Max)
	a.EqInt(45, face.Age.Min)
	a.EqFloat64(0.372036, face.Age.Score)
	a.EqInt64(75, face.FaceLocation.Height)
	a.EqInt64(256, face.FaceLocation.Left)
	a.EqInt64(93, face.FaceLocation.Top)
	a.EqInt64(67, face.FaceLocation.Width)
	a.EqStr("MALE", face.Gender.Gender)
	a.EqFloat64(0.99593, face.Gender.Score)
	a.EqStr("Barack Obama", face.Identity.Name)
	a.EqFloat64(0.989013, face.Identity.Score)
	a.EqStr("/people/politicians/democrats/barackobama", face.Identity.TypeHierarchy)
}
