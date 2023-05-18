package core

import (
	"errors"
	"fmt"
)

type ErrorCode uint8

const (
	ErrorUnknown ErrorCode = iota + 1
	ErrorConflict
	ErrorNotFound
	ErrorInvalid
	ErrorInternal
)

func (s ErrorCode) String() string {
	switch s {
	case ErrorConflict:
		return "CONFLICT"
	case ErrorNotFound:
		return "NOT_FOUND"
	case ErrorInvalid:
		return "INVALID"
	case ErrorInternal:
		return "INTERNAL"
	}
	panic(fmt.Sprintf("unknown job status %d", s))
}

// Error represents an application-specific error. Application errors can be
// unwrapped by the caller to extract out the code & message.
//
// Any non-application error (such as a disk error) should be reported as an
// EINTERNAL error and the human user should only see "Internal error" as the
// message. These low-level internal error details should only be logged and
// reported to the operator of the application (not the end user).
type Error struct {
	// Machine-readable error code.
	Code ErrorCode

	// Human-readable error message.
	Message string
}

// Implements the error interface. Not used by the application otherwise.
func (e *Error) Error() string {
	return fmt.Sprintf("{code: %s, message: %s}", e.Code.String(), e.Message)
}

// Unwraps an application error and returns its code.
func UnwrapCode(err error) ErrorCode {
	var e *Error
	if err == nil {
		return ErrorUnknown
	} else if errors.As(err, &e) {
		return e.Code
	}
	return ErrorInternal
}

// Unwraps an application error and returns its message.
func UnwrapMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error"
}

func Errorf(code ErrorCode, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
