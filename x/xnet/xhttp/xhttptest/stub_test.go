package xhttptest

import (
	"io/ioutil"
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

func TestStubFile(t *testing.T) {
	a := assert.New(t)
	client := StubFile(map[string]string{
		"http://www.example.com/": "./fixtures/stubfile.txt",
	}, &http.Client{})
	resp, err := client.Get("http://www.example.com/")
	a.Nil(err)
	body, err := ioutil.ReadAll(resp.Body)
	a.Nil(err)
	a.EqStr("OK", string(body))
}
