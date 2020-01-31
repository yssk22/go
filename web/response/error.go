package response

import (
	"fmt"
	"io"

	"context"
)

type _error struct {
	err error
}

func (e *_error) Render(ctx context.Context, w io.Writer) {
	fmt.Fprintf(w, "%v", e.err)
}

// NewError returns a text response
func NewError(ctx context.Context, e error) *Response {
	return NewResponseWithStatus(ctx, &_error{e}, HTTPStatusInternalServerError)
}

// NewErrorWithStatus returns an error response with the given status code
func NewErrorWithStatus(ctx context.Context, e error, code HTTPStatus) *Response {
	return NewResponseWithStatus(ctx, &_error{e}, code)
}
