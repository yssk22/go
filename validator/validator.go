// Package validator provides types and functions for validator implementations
package validator

import (
	"fmt"
	"reflect"
)

// Validatable is an interface for an object to be validated
type Validatable interface {
	Value() interface{}
	Get(string) (interface{}, error)
}

// Validator is a struct that contains validator functions
type Validator struct {
	fields map[string]*List // field validators
	*List                   // object validators
}

// NewValidator returns a new *BaseValidator
func NewValidator() *Validator {
	return &Validator{
		fields: make(map[string]*List),
		List:   NewList(),
	}
}

// Field returns a *List for the given field.
func (v *Validator) Field(name string) *List {
	_, ok := v.fields[name]
	if !ok {
		v.fields[name] = NewList()
	}
	return v.fields[name]
}

// Eval to evaluate the object t and return the validation result.
// nil would be returned on no validation errors.
func (v *Validator) Eval(i interface{}) *ValidationError {
	var validatable Validatable
	if v, ok := i.(Validatable); ok {
		validatable = v
	} else {
		validatable = NewValidatable(i)
	}
	ve := make(map[string][]*FieldError)
	for name, list := range v.fields {
		value, err := validatable.Get(name)
		if err != nil {
			ferr := &FieldError{
				Message: err.Error(),
			}
			if _, ok := ve[name]; ok {
				ve[name] = append(ve[name], ferr)
			} else {
				ve[name] = []*FieldError{ferr}
			}
			continue
		}
		for _, f := range list.funcs {
			err := f(value)
			if err != nil {
				if _, ok := ve[name]; ok {
					ve[name] = append(ve[name], err)
				} else {
					ve[name] = []*FieldError{err}
				}
			}
		}
	}

	// object validation
	const objectKey = ""
	for _, f := range v.funcs {
		err := f(validatable.Value())
		if err != nil {
			if _, ok := ve[objectKey]; ok {
				ve[objectKey] = append(ve[objectKey], err)
			} else {
				ve[objectKey] = []*FieldError{err}
			}
		}
	}
	if len(ve) > 0 {
		return &ValidationError{ve, validatable.Value()}
	}
	return nil
}

type validatable struct {
	v reflect.Value
	i interface{}
}

// NewValidatable returns a Validatable for any types of struct object
func NewValidatable(i interface{}) Validatable {
	v := reflect.ValueOf(i)
	k := v.Kind()
	if k == reflect.Ptr {
		v = reflect.Indirect(v)
		k = v.Kind()
	}
	if k != reflect.Struct {
		panic(fmt.Errorf("could not create a Validatable for non-struct type (%s)", v.Type().Name()))
	}
	return &validatable{
		v: v,
		i: i,
	}
}

func (v *validatable) Value() interface{} {
	return v.i
}

func (v *validatable) Get(name string) (interface{}, error) {
	val := v.v.FieldByName(name)
	if val.IsValid() {
		return val.Interface(), nil
	}
	return nil, nil
}
