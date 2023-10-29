package errors

import (
	"fmt"
)

type wrappedError struct {
	msg   string
	cause error
}

func (w *wrappedError) Error() string {
	return fmt.Sprintf("%s: %s", w.msg, w.cause)
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &wrappedError{
		cause: err,
		msg:   message,
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &wrappedError{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}
