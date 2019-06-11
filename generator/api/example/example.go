package example

import (
	"context"
	"errors"

	"github.com/yssk22/go/generator/api/example/types"
)

// Example struct
type Example struct {
	ID string `json:"id"`
}

// @api path=/path/to/example/:param/:param2/
func getExample(ctx context.Context, param string, param2 string) (*types.ResponseData, error) {
	a := &types.ResponseData{
		StrVal: param,
	}
	return a, nil
}

// @api path=/path/to/example/:param/2/
func getExampleWithExtraParam(ctx context.Context, param string, query *types.RequestParams) (*types.RequestParams, error) {
	return query, nil
}

// @api path=/path/to/example/:param/
func createExample(ctx context.Context, param string, e *Example) (string, error) {
	a := &types.ResponseData{}
	return a.StrVal, nil
}

// @api path=/path/to/example2/:param/ format=json
func createExample2(ctx context.Context, param string, e *types.ComplexRequestParams) (string, error) {
	a := &types.ResponseData{}
	return a.StrVal, nil
}

// @api path=/path/to/example/:param/
func updateExample(ctx context.Context, param string, e *Example) (*types.ResponseData, error) {
	a := &types.ResponseData{}
	return a, nil
}

// @api path=/path/to/example/:param/
func deleteExample(ctx context.Context, param string, query *types.RequestParams) (*types.ResponseData, error) {
	a := &types.ResponseData{}
	return a, nil
}

// @api path=/path/to/example/:param/:param2/always_ok/
func getExampleAlwaysOK(ctx context.Context, param string, param2 string) {
}

// @api path=/path/to/example/:param/:param2/only_error/
func getExampleOnlyError(ctx context.Context, param string, param2 string) error {
	return errors.New("OnlyError")
}
