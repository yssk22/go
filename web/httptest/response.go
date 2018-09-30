package httptest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/yssk22/go/web/response"
	"github.com/yssk22/go/x/xtesting/assert"
)

// Assert is a wrapper for github.com/yssk22/go/x/xtesting/assert.Assert and provides
// http specific assertions
type Assert struct {
	*assert.Assert
}

// NewAssert returns *Assert
func NewAssert(t *testing.T) *Assert {
	return &Assert{
		assert.New(t),
	}
}

// Status asserts the http status code
func (a *Assert) Status(expected response.HTTPStatus, res *httptest.ResponseRecorder, msgContext ...interface{}) {
	if expected != response.HTTPStatus(res.Code) {
		if len(msgContext) > 0 {
			a.Failure(expected, res.Code, msgContext...)
		} else {
			a.Failure(expected, res.Code, "**** HTTP Body ****\n\t%s", res.Body)
		}
	}
}

// Header asserts the header value
func (a *Assert) Header(expected string, res *httptest.ResponseRecorder, fieldName string, msgContext ...interface{}) {
	a.EqStr(expected, res.Header().Get(fieldName), msgContext...)
}

// Body asserts the body string
func (a *Assert) Body(expected string, res *httptest.ResponseRecorder, msgContext ...interface{}) {
	a.EqStr(expected, res.Body.String(), msgContext...)
}

// Cookie asserts the cookie name and extract it as *http.Cookie
func (a *Assert) Cookie(res *httptest.ResponseRecorder, name string, msgContext ...interface{}) *http.Cookie {
	rawCookies, ok := res.Header()["Set-Cookie"]
	if !ok {
		a.Failure("Set-Cookie header exists", nil, msgContext)
	}
	var buff []string
	buff = append(buff, "GET / HTTP/1.1")
	for _, v := range rawCookies {
		buff = append(buff, fmt.Sprintf("Cookie: %s", v))
	}
	buff = append(buff, "\r\n")
	rawRequest := strings.Join(buff, "\r\n")
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))
	if err != nil {
		panic("err")
	}
	c, err := req.Cookie(name)
	if err != nil {
		a.Failure("No cookie %s presents", req.Cookies, msgContext)
	}
	return c
}

// JSON asserts the body string as as json and returns the result as interface{}
func (a *Assert) JSON(v interface{}, res *httptest.ResponseRecorder, msgContext ...interface{}) {
	var body = res.Body.Bytes()
	err := json.Unmarshal(body, v)
	if err != nil {
		a.Failure(
			fmt.Sprintf("%v of %s", err, reflect.TypeOf(v).Name()),
			string(body),
			msgContext...,
		)
	}
}
