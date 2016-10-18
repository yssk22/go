package errors

import "fmt"

// F returns a short hand for fmt.Errorf
func F(format string, v ...interface{}) error {
	return fmt.Errorf(format, v...)
}
