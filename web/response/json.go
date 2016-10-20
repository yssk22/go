package response

import (
	"encoding/json"
	"net/http"
)

// JSON implements application/json resposne.
type JSON struct {
	Code   HTTPStatus
	Header *Header
	Data   interface{}
}

// NewJSON returns *JSON reponse
func NewJSON(v interface{}) *JSON {
	return NewJSONWithCode(v, HTTPStatusOK)
}

const contentTypeJSON = "application/json; charset=utf-8"

// NewJSONWithCode returns *JSON reponse with the given status code
func NewJSONWithCode(v interface{}, code HTTPStatus) *JSON {
	header := NewHeader()
	header.Set(contentTypeKey, contentTypeJSON)
	return &JSON{
		Code:   code,
		Header: header,
		Data:   v,
	}
}

// Render writes text response
func (j *JSON) Render(w http.ResponseWriter) bool {
	j.Header.Render(w)
	w.WriteHeader(int(j.Code))
	if err := json.NewEncoder(w).Encode(j.Data); err != nil {
		panic(err)
	}
	return true
}
