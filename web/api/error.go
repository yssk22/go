package api

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/yssk22/go/web/response"
	"github.com/yssk22/go/x/xerrors"
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
func (err *Error) ToResponse(ctx context.Context) *response.Response {
	if err.Status == response.HTTPStatus(0) {
		return response.NewJSONWithStatus(ctx, err, response.HTTPStatusBadRequest)
	}
	return response.NewJSONWithStatus(ctx, err, err.Status)
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%s]%s: %s", err.Status, err.Code, err.Message)
}

// NewErrorResponse returns response.Response with the complient format
func NewErrorResponse(ctx context.Context, e error) *response.Response {
	apie, ok := e.(*Error)
	if ok {
		return apie.ToResponse(ctx)
	}
	stacks := xerrors.Unwrap(e)
	log.Println(fmt.Sprintf("internal error occurred - %s", strings.Join(stacks, "\n")))
	return (&Error{
		Code:    "internal_server_error",
		Message: "unexpected server error occurred. please try again later.",
		Status:  response.HTTPStatusInternalServerError,
	}).ToResponse(ctx)
}
