package response

// NewRedirect returns *Responose for redirect
func NewRedirect(url string) *Response {
	return NewRedirectWithStatus(url, HTTPStatusSeeOther)
}

// NewRedirectWithStatus returns *Responose for redirect
func NewRedirectWithStatus(url string, status HTTPStatus) *Response {
	switch status {
	case HTTPStatusMovedParmanently, HTTPStatusFound, HTTPStatusSeeOther:
		res := NewResponseWithStatus(NoContent, status)
		res.Header.Add("Location", url)
		return res
	default:
		panic("Invalid status code for redirect")
	}
}
