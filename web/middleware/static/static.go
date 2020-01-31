package static

import (
	"io/ioutil"
	"os"

	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
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
		return response.NewError(req.Context(), err)
	}
	return response.NewText(req.Context(), string(buff))
}
