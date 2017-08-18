package xerrors

import "fmt"

type wrapped struct {
	cause   error
	message string
}

func (w *wrapped) Error() string {
	return fmt.Sprintf("%s: %s", w.message, w.cause)
}

// Wrap wraps an error with an additional message.
func Wrap(err error, format string, args ...interface{}) error {
	return &wrapped{
		cause:   err,
		message: fmt.Sprintf(format, args...),
	}
}
