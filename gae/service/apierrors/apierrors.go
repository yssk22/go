package apierrors

import (
	"github.com/speedland/go/web/response"
)

// Error represents API error
type Error struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Status  response.HTTPStatus `json:"-"`
}

// ToResponse returns *response.Response object for this error
func (err *Error) ToResponse() *response.Response {
	if err.Status == response.HTTPStatus(0) {
		return response.NewJSONWithStatus(err, response.HTTPStatusBadRequest)
	}
	return response.NewJSONWithStatus(err, err.Status)
}

// General API errors
var (
	Forbidden = &Error{
		Code:    "forbidden",
		Message: "you do not have an access to the resource",
		Status:  response.HTTPStatusForbidden,
	}
)
