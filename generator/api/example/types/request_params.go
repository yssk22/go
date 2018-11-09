package types

import (
	"context"

	"github.com/yssk22/go/validator"
	"github.com/yssk22/go/web/api"
)

// RequestParams is an example of typed request parameters
type RequestParams struct {
	StrVal         string  `json:"str_val"`
	StrValDefault  string  `json:"str_val_default" default:"foo"`
	StrValRequired string  `json:"str_val_required" validate:"required"`
	StrPtr         *string `json:"str_ptr"`
	StrPtrDefault  *string `json:"str_ptr_default" default:"bar"`
	StrPtrRequired *string `json:"str_ptr_required" validate:"required"`

	IntVal int `json:"int_val"`
}

// Validate validates the parameter
func (r *RequestParams) Validate(ctx context.Context, errors api.FieldErrorCollection) error {
	errors.Add("int_val", validator.Int().Min(0).Max(10).Validate(r.IntVal))
	return nil
}

// ComplexRequestParams is a struct to represent complex types
type ComplexRequestParams struct {
	Str      string         `json:"str"`
	Inner    *RequestParams `json:"inner"`
	IntArray []int          `json:"int_array"`
	MyEnum   MyEnum         `json:"my_enum"`
}
