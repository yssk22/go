package example

import (
	"time"

	"github.com/yssk22/go/types"
)

// Example to test struct
// @flow
type Example struct {
	BoolVal   bool       `json:"bool_val"`
	BoolPtr   *bool      `json:"bool_ptr"`
	IntVaL    int        `json:"int_val"`
	IntPtr    *int       `json:"int_ptr"`
	FloatVal  float64    `json:"float_val"`
	FloatPtr  *float64   `json:"float_ptr"`
	StringVal string     `json:"string_val"`
	StringPtr *string    `json:"string_ptr"`
	TimeVal   time.Time  `json:"time_val"`
	TimePtr   *time.Time `json:"time_ptr"`
	InnerVal  Inner      `json:"inner_val"`
	InnerPtr  *Inner     `json:"inner_ptr"`
	Imported  types.RGB  `json:"rgb"`
}

// Inner to test inner object
// @flow
type Inner struct {
	BoolVal   bool       `json:"bool_val"`
	BoolPtr   *bool      `json:"bool_ptr"`
	IntVaL    int        `json:"int_val"`
	IntPtr    *int       `json:"int_ptr"`
	FloatVal  float64    `json:"float_val"`
	FloatPtr  *float64   `json:"float_ptr"`
	StringVal string     `json:"string_val"`
	StringPtr *string    `json:"string_ptr"`
	TimeVal   time.Time  `json:"time_val"`
	TimePtr   *time.Time `json:"time_ptr"`
}
