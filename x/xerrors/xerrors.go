package xerrors

import "fmt"

// F returns a short hand for fmt.Errorf
func F(format string, v ...interface{}) error {
	return fmt.Errorf(format, v...)
}

// MustNil checks err is nil, or panic.
func MustNil(err error) {
	if err != nil {
		panic(err)
	}
}

// MustValue returns v if an err is nil, or panic.
func MustValue(v interface{}, err error) interface{} {
	MustNil(err)
	return v
}

// MultiError is an error collection as a single error.
// error[i] might be nil if there is no error.
type MultiError []error

// NewMultiError creates MultiError instance with the given size.
func NewMultiError(size int) MultiError {
	return MultiError(make([]error, size))
}

// Error implemnts error.Error()
func (me MultiError) Error() string {
	var firstError error
	var errorCount int
	for _, e := range me {
		if e != nil {
			if firstError == nil {
				firstError = e
			}
			errorCount++
		}
	}
	switch errorCount {
	case 0:
		return "no error"
	case 1:
		return firstError.Error()
	}
	return fmt.Sprintf("%s (and %d other errors)", firstError.Error(), errorCount-1)
}

// HasError returns if there is an error in the errors.
func (me MultiError) HasError() bool {
	for _, e := range me {
		if e != nil {
			return true
		}
	}
	return false
}

// ToReturn returns itself if MultiError has an error, otherwise nil
func (me MultiError) ToReturn() error {
	if me.HasError() {
		return me
	}
	return nil
}
