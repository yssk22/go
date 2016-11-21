package visualrecognition

import (
	"net/http"
	"os"
	"testing"
	"x/assert"

	"github.com/speedland/go/x/xnet/xhttp/xhttptest"
	"golang.org/x/net/context"
)

func TestClient_DetectByURL(t *testing.T) {
	a := assert.New(t)
	stubClient := xhttptest.StubFile(
		map[string]string{
			"https://gateway-a.watsonplatform.net/visual-recognition/api/v3/detect_faces?api_key=abc&url=https%3A%2F%2Fgithub.com%2Fwatson-developer-cloud%2Fdoc-tutorial-downloads%2Fraw%2Fmaster%2Fvisual-recognition%2Fprez.jpg&version=2016-05-20": "./fixtures/fd.json",
		},
		&http.Client{},
	)
	c := NewClient("abc", stubClient)
	resp, err := c.DetectFacesByURL(context.Background(), "https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/prez.jpg")
	a.Nil(err)
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

func TestClient_DetectByURL_real_access(t *testing.T) {
	a := assert.New(t)
	apiKey := os.Getenv("WATSON_VISUAL_RECOGNITION_API_KEY")
	if apiKey == "" {
		t.Skipf("Environment key %q is not set.", "WATSON_VISUAL_RECOGNITION_API_KEY")
	}
	c := NewClient(apiKey, http.DefaultClient)
	resp, err := c.DetectFacesByURL(context.Background(), "https://github.com/watson-developer-cloud/doc-tutorial-downloads/raw/master/visual-recognition/prez.jpg")
	a.Nil(err)
	a.EqInt(1, len(resp.Images))
	a.EqInt(1, len(resp.Images[0].Faces))
}
