package response

import "net/http"

// Test implements text/plain resposne.
type Text struct {
	Code    HTTPStatus
	Header  *Header
	Content string
}

// NewText returns *Text reponse
func NewText(content string) *Text {
	return NewTextWithCode(content, HTTPStatusOK)
}

// NewTextWithCode returns *Text reponse with the given status code
func NewTextWithCode(content string, code HTTPStatus) *Text {
	header := NewHeader()
	header.Set("content-type", "text/plain; charset=utf-9")
	return &Text{
		Code:    code,
		Header:  header,
		Content: content,
	}
}

// Render writes text response
func (t *Text) Render(w http.ResponseWriter) bool {
	t.Header.Render(w)
	w.WriteHeader(int(t.Code))
	w.Write([]byte(t.Content))
	return true
}
