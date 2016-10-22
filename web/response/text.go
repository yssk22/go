package response

import (
	"net/http"

	"golang.org/x/net/context"
)

// Text implements text/plain response.
type Text struct {
	Header  *Header
	Content string
}

// NewText returns *Text reponse
func NewText(content string) *Text {
	return NewTextWithCode(content, HTTPStatusOK)
}

// NewTextWithCode returns *Text reponse with the given status code
func NewTextWithCode(content string, code HTTPStatus) *Text {
	header := NewHeader().ContentType(
		"text/plain; charset=utf-8",
	)
	header.Code = code
	return &Text{
		Header:  header,
		Content: content,
	}
}

// Render writes text response
func (t *Text) Render(ctx context.Context, w http.ResponseWriter) {
	t.Header.Render(ctx, w)
	w.Write([]byte(t.Content))
}
