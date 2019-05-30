package api

import (
	"context"
	"fmt"
)

// Validatable is an interface to run validate method for the parameter struct after parsing.
type Validatable interface {
	Validate(ctx context.Context, errors FieldErrorCollection) error
}

// IntValidator validates int value
type IntValidator struct {
	funcs [](func(int) error)
}

// Min validates the minimum
func (v *IntValidator) Min(i int) *IntValidator {
	v.funcs = append(v.funcs, func(value int) error {
		if value < i {
			return fmt.Errorf("must be more than or equal to %d", i)
		}
		return nil
	})
	return v
}

// Validate validates the int value
func (v *IntValidator) Validate(i int) error {
	for _, v := range v.funcs {
		if err := v(i); err != nil {
			return err
		}
	}
	return nil
}
