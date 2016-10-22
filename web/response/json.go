package response

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

// JSON implements application/json resposne.
type JSON struct {
	Header *Header
	Data   interface{}
}

// NewJSON returns *JSON reponse
func NewJSON(v interface{}) *JSON {
	return NewJSONWithCode(v, HTTPStatusOK)
}

// NewJSONWithCode returns *JSON reponse with the given status code
func NewJSONWithCode(v interface{}, code HTTPStatus) *JSON {
	header := NewHeader().ContentType(
		"application/json; charset=utf-8",
	)
	header.Code = code
	return &JSON{
		Header: header,
		Data:   v,
	}
}

// Render writes text response
func (j *JSON) Render(ctx context.Context, w http.ResponseWriter) {
	j.Header.Render(ctx, w)
	if err := json.NewEncoder(w).Encode(j.Data); err != nil {
		panic(err)
	}
}
