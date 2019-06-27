package validator

import "fmt"

// String returns a StringValidator
func String() StringValidator {
	return StringValidator(make([](func(string) error), 0))
}

// StringValidator to validate the string value
type StringValidator []func(string) error

// Validate validates the given value
func (funcs StringValidator) Validate(s string) error {
	for _, v := range funcs {
		if err := v(s); err != nil {
			return err
		}
	}
	return nil
}

// NonEmpty validates s is not an empty string
func (funcs StringValidator) NonEmpty() StringValidator {
	funcs = append(funcs, func(v string) error {
		if len(v) == 0 {
			return fmt.Errorf("must not be empty")
		}
		return nil
	})
	return funcs
}

// Min validates the minimum length of string
func (funcs StringValidator) Min(i int) StringValidator {
	funcs = append(funcs, func(v string) error {
		if len(v) < i {
			return fmt.Errorf("must be more than or equal to %d", i)
		}
		return nil
	})
	return funcs
}

// Max validates the max length of string
func (funcs StringValidator) Max(i int) StringValidator {
	funcs = append(funcs, func(v string) error {
		if len(v) > i {
			return fmt.Errorf("must be less than or equal to %d", i)
		}
		return nil
	})
	return funcs
}

// Range validates the maximum
func (funcs StringValidator) Range(s int, t int) StringValidator {
	funcs = append(funcs, func(v string) error {
		l := len(v)
		if l < s || l > t {
			return fmt.Errorf("must be [%d-%d] charactors", s, t)
		}
		return nil
	})
	return funcs
}
