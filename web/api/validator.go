package api

import (
	"context"
)

// Validatable is an interface to run validate method for the parameter struct after parsing.
type Validatable interface {
	Validate(ctx context.Context, errors FieldErrorCollection) error
}

// RunValidation runs a v#Validate to collect FieldErrorCollection and/or error
func RunValidation(ctx context.Context, v Validatable) (FieldErrorCollection, error) {
	fieldErrors := NewFieldErrorCollection()
	err := v.Validate(ctx, fieldErrors)
	return fieldErrors, err
}