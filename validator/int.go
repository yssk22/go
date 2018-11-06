package validator

import "fmt"

// Int returns a IntValidator
func Int() IntValidator {
	return IntValidator(make([](func(int) error), 0))
}

// IntValidator to validate the int value
type IntValidator []func(int) error

// Validate validates the given value
func (funcs IntValidator) Validate(i int) error {
	for _, v := range funcs {
		if err := v(i); err != nil {
			return err
		}
	}
	return nil
}

// Min validates the minimum
func (funcs IntValidator) Min(i int) IntValidator {
	funcs = append(funcs, func(v int) error {
		if v < i {
			return fmt.Errorf("must be more than or equal to %d", i)
		}
		return nil
	})
	return funcs
}

// Max validates the maximum
func (funcs IntValidator) Max(i int) IntValidator {
	funcs = append(funcs, func(v int) error {
		if v > i {
			return fmt.Errorf("must be less than or equal to %d", i)
		}
		return nil
	})
	return funcs
}

// Range validates the maximum
func (funcs IntValidator) Range(s int, t int) IntValidator {
	funcs = append(funcs, func(v int) error {
		if v < s || v > t {
			return fmt.Errorf("must be %d ~ %d", s, t)
		}
		return nil
	})
	return funcs
}
