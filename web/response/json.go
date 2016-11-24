package response

import (
	"encoding/json"
	"io"

	"golang.org/x/net/context"
)

// UseFormattedJSON is a configuration variable about if json object is formatted or not.
var UseFormattedJSON = false

type _json struct {
	data interface{}
}

func (j _json) Render(ctx context.Context, w io.Writer) {
	if UseFormattedJSON {
		buff, err := json.MarshalIndent(j.data, "", "    ")
		if err != nil {
			panic(err)
		}
		w.Write(buff)
		return
	}
	if err := json.NewEncoder(w).Encode(j.data); err != nil {
		panic(err)
	}
}

// NewJSON returns a JSON response
func NewJSON(v interface{}) *Response {
	return NewJSONWithStatus(v, HTTPStatusOK)
}

// NewJSONWithStatus returns a JSON formatted response with the given status code
func NewJSONWithStatus(v interface{}, code HTTPStatus) *Response {
	res := NewResponseWithStatus(&_json{v}, code)
	res.Header.Set(ContentType, "application/json; charset=utf-8")
	return res
}
