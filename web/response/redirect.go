package response

import "context"

// NewRedirect returns *Responose for redirect
func NewRedirect(ctx context.Context, url string) *Response {
	return NewRedirectWithStatus(ctx, url, HTTPStatusSeeOther)
}

// NewRedirectWithStatus returns *Responose for redirect
func NewRedirectWithStatus(ctx context.Context, url string, status HTTPStatus) *Response {
	switch status {
	case HTTPStatusMovedParmanently, HTTPStatusFound, HTTPStatusSeeOther:
		res := NewResponseWithStatus(ctx, NoContent, status)
		res.Header.Set("Location", url)
		return res
	default:
		panic("Invalid status code for redirect")
	}
}
