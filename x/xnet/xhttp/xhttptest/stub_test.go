package xhttptest

import (
	"net/http"
	"testing"
	"x/assert"
)

func TestStub(t *testing.T) {
	a := assert.New(t)
	client := Stub(nil, &http.Client{})
	_, err := client.Get("http://www.example.com/")
	a.NotNil(err)
	a.EqStr("Get http://www.example.com/: forbitten by xhttptest.Rewriter", err.Error())
}
