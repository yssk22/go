package api

import (
	"fmt"

	"github.com/yssk22/go/web/response"
)

// Error is an object to represent API error
type Error struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Trace   []string            `json:"trace,omitempty"`
	Extra   interface{}         `json:"extra"`
	Status  response.HTTPStatus `json:"-"`
}

// ToResponse returns *response.Response object for this error
func (err *Error) ToResponse() *response.Response {
	if err.Status == response.HTTPStatus(0) {
		return response.NewJSONWithStatus(err, response.HTTPStatusBadRequest)
	}
	return response.NewJSONWithStatus(err, err.Status)
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%s]%s: %s", err.Status, err.Code, err.Message)
}

// NewErrorResponse returns response.Response with the complient format
func NewErrorResponse(e error) *response.Response {
	apie, ok := e.(*Error)
	if ok {
		return apie.ToResponse()
	}
	// TODO: capture stack here
	return (&Error{
		Code:    "internal_server_error",
		Message: "unexpected server error occurred. please try again later.",
		Status:  response.HTTPStatusInternalServerError,
	}).ToResponse()
}

// General API errors
var (
	BadRequest = &Error{
		Code:    "bad_request",
		Message: "we cannot process your request",
		Status:  response.HTTPStatusForbidden,
	}
	Forbidden = &Error{
		Code:    "forbidden",
		Message: "you do not have an access to the resource",
		Status:  response.HTTPStatusForbidden,
	}
	NotFound = &Error{
		Code:    "not_found",
		Message: "the requested resource is not found on the server",
		Status:  response.HTTPStatusNotFound,
	}
	ServerError = &Error{
		Code:    "internal_server_error",
		Message: "unexpected server error occurred. please try again later.",
		Status:  response.HTTPStatusInternalServerError,
	}
)