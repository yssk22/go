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
