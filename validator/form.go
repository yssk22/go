package validator

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/yssk22/go/x/xtime"
)

type fieldType int

// Available FieldType values
const (
	fieldTypeString fieldType = iota
	fieldTypeBool
	fieldTypeInt
	fieldTypeFloat
	fieldTypeDate
	fieldTypeTime
	fieldTypeDateTime
	fieldTypeDuration
)

// FormValidator is a validator to validate url.Value
type FormValidator struct {
	*Validator
	fieldTypes map[string]fieldType
}

// NewFormValidator returns a new *FormValidator
func NewFormValidator() *FormValidator {
	return &FormValidator{
		Validator:  NewValidator(),
		fieldTypes: make(map[string]fieldType),
	}
}

// BoolField - a named field for the bool value.
func (f *FormValidator) BoolField(name string) *List {
	f.fieldTypes[name] = fieldTypeBool
	return f.Field(name)
}

// IntField - a named field for the int value.
func (f *FormValidator) IntField(name string) *List {
	f.fieldTypes[name] = fieldTypeInt
	return f.Field(name)
}

// FloatField - a named field for the float value.
func (f *FormValidator) FloatField(name string) *List {
	f.fieldTypes[name] = fieldTypeFloat
	return f.Field(name)
}

// DateField - a named field for the time.Time value (formatted as xtime.ParseDate compatible).
func (f *FormValidator) DateField(name string) *List {
	f.fieldTypes[name] = fieldTypeDate
	return f.Field(name)
}

// TimeField - a named field for the time.Time value (formatted as xtime.ParseTime compatible).
func (f *FormValidator) TimeField(name string) *List {
	f.fieldTypes[name] = fieldTypeTime
	return f.Field(name)
}

// DateTimeField - a named field for the time.Time value (formatted as xtime.Parse compatible).
func (f *FormValidator) DateTimeField(name string) *List {
	f.fieldTypes[name] = fieldTypeDateTime
	return f.Field(name)
}

// DurationField - a named field for the time.Duration value.
func (f *FormValidator) DurationField(name string) *List {
	f.fieldTypes[name] = fieldTypeDuration
	return f.Field(name)
}

// Eval evaluates the `values` and returns the error if the validation fails.
func (f *FormValidator) Eval(values url.Values) *ValidationError {
	return f.Validator.Eval(
		&formValidatable{
			values:     values,
			fieldTypes: f.fieldTypes,
		},
	)
}

type formValidatable struct {
	values     url.Values
	fieldTypes map[string]fieldType
}

func (fv *formValidatable) Value() interface{} {
	return fv.values
}

func (fv *formValidatable) Get(field string) (interface{}, error) {
	values := fv.values[field]
	typed := fv.fieldTypes[field]
	if values != nil {
		var v interface{}
		var e error
		if len(values) > 0 {
			switch typed {
			case fieldTypeBool:
				v, e = values[0] == "true" || values[0] == "1", nil
			case fieldTypeInt:
				v, e = strconv.ParseInt(values[0], 10, 64)
			case fieldTypeFloat:
				v, e = strconv.ParseFloat(values[0], 64)
			case fieldTypeDate:
				v, e = xtime.ParseDate(values[0], time.UTC, 0)
			case fieldTypeTime:
				v, e = xtime.ParseTime(values[0], time.UTC)
			case fieldTypeDateTime:
				v, e = xtime.Parse(values[0])
			case fieldTypeDuration:
				v, e = time.ParseDuration(values[0])
			default: // string
				v, e = values[0], nil
			}
		}
		if e != nil {
			return nil, fmt.Errorf("Invalid format %q", values[0])
		}
		return v, nil
	}
	return nil, nil
}
