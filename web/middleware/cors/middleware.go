package cors

import (
	"log"

	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
)

type Middleware struct {
	AllowOrigin string
}

func NewMiddleware(allowOrigin string) *Middleware {
	return &Middleware{
		AllowOrigin: allowOrigin,
	}
}

func (m *Middleware) Process(req *web.Request, next web.NextHandler) *response.Response {
	log.Println("New request")
	resp := next(req)
	resp.Header.Add("Access-Control-Allow-Origin", m.AllowOrigin)
	return resp
}
