package static

import (
	"io/ioutil"
	"os"

	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
)

// ServeFile returns web.Handler that serve a file content on `path`
func ServeFile(path string) web.Handler {
	return &serveFile{
		path,
	}
}

type serveFile struct {
	path string
}

func (s *serveFile) Process(req *web.Request, next web.NextHandler) *response.Response {
	f, err := os.Open(s.path)
	if err != nil {
		return nil
	}
	defer f.Close()
	buff, err := ioutil.ReadAll(f)
	if err != nil {
		return response.NewError(err)
	}
	return response.NewText(string(buff))
}
