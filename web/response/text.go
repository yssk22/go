package response

import (
	"fmt"
	"io"

	"golang.org/x/net/context"
)

type _text struct {
	v interface{}
}

func (t *_text) Render(ctx context.Context, w io.Writer) {
	fmt.Fprintf(w, "%v", t.v)
}

// NewText returns a text response
func NewText(s interface{}) *Response {
	return NewTextWithStatus(s, HTTPStatusOK)
}

// NewTextWithStatus returns a text formatted response with the given status code
func NewTextWithStatus(s interface{}, code HTTPStatus) *Response {
	res := NewResponseWithStatus(&_text{s}, code)
	res.Header.Set(ContentType, "plain/text; charset=utf-8")
	return res
}
