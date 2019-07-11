package api

import "github.com/yssk22/go/web/response"

var _ok = map[string]bool{
	"OK": true,
}

// OK returns { "OK": true }
func OK() *response.Response {
	return response.NewJSON(_ok)
}

// General API errors
var (
	BadRequest = &Error{
		Code:    "bad_request",
		Message: "we cannot process your request",
		Status:  response.HTTPStatusBadRequest,
	}
	Unauthorized = &Error{
		Code:    "unauthorized",
		Message: "authentication is required to proceed your request",
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
