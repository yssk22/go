package xhttptest

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// GetCookie returns *http.Cookie from httptes.ResponseRecorder
func GetCookie(w *httptest.ResponseRecorder, name string) (*http.Cookie, error) {
	rawCookies, ok := w.Header()["Set-Cookie"]
	if !ok {
		return nil, fmt.Errorf("'Set-Cookie' header does not present")
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
		return nil, err
	}
	c, err := req.Cookie(name)
	if err != nil {
		return nil, fmt.Errorf("Cookie %q not found")
	}
	return c, nil
}
