package types

import (
	"time"
)

// ResponseData is an example of typed response parameters
type ResponseData struct {
	StrVal string  `json:"str_val"`
	StrPtr *string `json:"str_ptr"`

	IntVal int  `json:"int_val"`
	IntPtr *int `json:"int_ptr"`

	FloatVal float64  `json:"float_val"`
	FloatPtr *float64 `json:"float_ptr"`

	BoolVal bool  `json:"bool_val"`
	BoolPtr *bool `json:"bool_ptr"`

	TimeVal time.Time  `json:"time_val"`
	TimePtr *time.Time `json:"time_ptr"`

	Inner *InnerResponse `json:"inner"`
}

// InnerResponse is one of inner types of ResponseData
type InnerResponse struct {
	StrVal string `json:"str_val"`
}
