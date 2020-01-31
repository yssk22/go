package response

import (
	"context"
	"io/ioutil"
)

// NewStaticFile returns a static file response
func NewStaticFile(ctx context.Context, path string, contentType string) *Response {
	return NewStaticFileWithStatus(ctx, path, contentType, HTTPStatusOK)
}

// NewStaticFileWithStatus returns a text formatted response with the given status code
func NewStaticFileWithStatus(ctx context.Context, path string, contentType string, code HTTPStatus) *Response {
	buff, err := ioutil.ReadFile(path)
	var s string
	if err != nil {
		s = err.Error()
	} else {
		s = string(buff)
	}
	res := NewResponseWithStatus(ctx, &_text{s}, code)
	res.Header.Set(ContentType, contentType)
	return res
}
