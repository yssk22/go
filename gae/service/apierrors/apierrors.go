package apierrors

import (
	"context"
	"github.com/yssk22/go/web/response"
)

// Error represents API error
type Error struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Status  response.HTTPStatus `json:"-"`
}

// ToResponse returns *response.Response object for this error
func (err *Error) ToResponse(ctx context.Context) *response.Response {
	if err.Status == response.HTTPStatus(0) {
		return response.NewJSONWithStatus(ctx, err, response.HTTPStatusBadRequest)
	}
	return response.NewJSONWithStatus(ctx, err, err.Status)
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
