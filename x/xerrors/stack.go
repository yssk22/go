package xerrors

import (
	"fmt"

	"github.com/yssk22/go/x/xruntime"
)

type wrapped struct {
	cause   error
	message string
	frame   *xruntime.Frame
}

func (w *wrapped) Error() string {
	return fmt.Sprintf("%s: %s", w.message, w.cause)
}

// Wrap wraps an error with an additional message.
func Wrap(err error, format string, args ...interface{}) error {
	return &wrapped{
		cause:   err,
		message: fmt.Sprintf(format, args...),
		frame:   xruntime.CaptureCaller(),
	}
}

// Unwrap extract all error messages in the wrapped error
func Unwrap(err error) []string {
	var w *wrapped
	var ok bool
	var stack []string
	w, ok = err.(*wrapped)
	if !ok {
		return []string{err.Error()}
	}
	for ok {
		stack = append(stack, fmt.Sprintf("%s:%d - %s", w.frame.FullFilePath, w.frame.LineNumber, w.message))
		err = w.cause
		if err == nil {
			break
		}
		w, ok = err.(*wrapped)
		if !ok {
			stack = append(stack, fmt.Sprintf("unknown - %s", err))
			break
		}
	}
	return stack
}
