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
func NewError(e error) *Response {
	return NewResponseWithStatus(&_error{e}, HTTPStatusInternalServerError)
}

// NewErrorWithStatus returns an error response with the given status code
func NewErrorWithStatus(e error, code HTTPStatus) *Response {
	return NewResponseWithStatus(&_error{e}, code)
}
