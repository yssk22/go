package types

import (
	"time"
)

// RequestParams is an example of typed request parameters
type RequestParams struct {
	StrVal            string  `json:"str_val"`
	StrPtr            *string `json:"str_ptr"`
	StrPtrWithDefault string  `json:"str_ptr_default" default:"foo"`

	IntVal            int  `json:"int_val"`
	IntPtr            *int `json:"int_ptr"`
	IntPtrWithDefault int  `json:"int_ptr_default" default:"10"`

	FloatVal            float64  `json:"float_val"`
	FloatPtr            *float64 `json:"float_ptr"`
	FloatPtrWithDefault float64  `json:"float_ptr_default" default:"2.0"`

	BoolVal            bool  `json:"bool_val"`
	BoolPtr            *bool `json:"bool_ptr"`
	BoolPtrWithDefault bool  `json:"bool_ptr_default" default:"false"`

	TimeVal            time.Time  `json:"time_val"`
	TimePtr            *time.Time `json:"time_ptr"`
	TimePtrWithDefault time.Time  `json:"time_ptr_default" default:"false"`
}

// ComplexRequestParams is a struct to represent complex types
type ComplexRequestParams struct {
	Str string `json:"str"`
}
