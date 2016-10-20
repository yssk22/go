package response

import "net/http"

// Header is a Response implementation that sets header values
type Header struct {
	http.Header
}

// NewHeader returns a new Header
func NewHeader() *Header {
	return &Header{
		http.Header(
			make(map[string][]string),
		),
	}
}

// Render sets the response header values
func (h *Header) Render(w http.ResponseWriter) bool {
	header := w.Header()
	for key, v := range h.Header {
		for _, vv := range v {
			header.Add(key, vv)
		}
	}
	return false
}
