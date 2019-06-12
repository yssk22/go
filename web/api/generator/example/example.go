package example

import (
	"context"
	"errors"
	"fmt"

	"github.com/yssk22/go/web/api/generator/example/types"
)

// @api path=/path/to/example/:param/:param2/
func getExample(ctx context.Context, param string, param2 string) (*types.ResponseData, error) {
	a := &types.ResponseData{
		StrVal: param,
	}
	return a, nil
}

// @api path=/path/to/example/:param/2/
func getExampleWithExtraParam(ctx context.Context, param string, query *types.QueryParams) (*types.QueryParams, error) {
	return query, nil
}

// @api path=/path/to/example/:param/
func createExample(ctx context.Context, param string, e *types.BodyParams) (string, error) {
	a := &types.ResponseData{}
	return a.StrVal, nil
}

// @api path=/path/to/example2/:param/ format=json
func createExample2(ctx context.Context, param string, e *types.ComplexQueryParams) (string, error) {
	a := &types.ResponseData{}
	return a.StrVal, nil
}

// @api path=/path/to/example/:param/
func updateExample(ctx context.Context, param string, e *types.BodyParams) (*types.ResponseData, error) {
	a := &types.ResponseData{}
	return a, nil
}

// @api path=/path/to/example/:param/
func deleteExample(ctx context.Context, param string, query *types.QueryParams) (*types.ResponseData, error) {
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

type StructExample struct {
	response string
}

// @api path=/path/to/struct_example/:param/
func (se *StructExample) getStructExample(ctx context.Context, param string) (string, error) {
	return fmt.Sprintf("%s %s", se.response, param), nil
}
