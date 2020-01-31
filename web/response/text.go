package response

import (
	"fmt"
	"io"

	"context"
)

type _text struct {
	v interface{}
}

func (t *_text) Render(ctx context.Context, w io.Writer) {
	fmt.Fprintf(w, "%v", t.v)
}

// NewText returns a text response
func NewText(ctx context.Context, s interface{}) *Response {
	return NewTextWithStatus(ctx, s, HTTPStatusOK)
}

// NewTextWithStatus returns a text formatted response with the given status code
func NewTextWithStatus(ctx context.Context, s interface{}, code HTTPStatus) *Response {
	res := NewResponseWithStatus(ctx, &_text{s}, code)
	res.Header.Set(ContentType, "text/plain; charset=utf-8")
	return res
}
