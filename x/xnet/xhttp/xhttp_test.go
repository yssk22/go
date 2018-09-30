package xhttp

import (
	"net/http"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestAbsoluteURL(t *testing.T) {
	a := assert.New(t)
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	a.EqStr("http://localhost/bar", AbsoluteURL(req, "/bar"))
}
