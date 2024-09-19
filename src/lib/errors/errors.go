package errors

import (
	"fmt"
	"strings"
)

// Error represents a custom error type
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewError creates a new error instance
func NewError(code int, message string, details ...string) *Error {
	e := &Error{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		e.Details = strings.Join(details, ", ")
	}
	return e
}

// Error implements the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// ErrorCode returns the error code
func (e *Error) ErrorCode() int {
	return e.Code
}

// ErrorMessage returns the error message
func (e *Error) ErrorMessage() string {
	return e.Message
}

// ErrorDetails returns the error details
func (e *Error) ErrorDetails() string {
	return e.Details
}

// IsError checks if the error is of type Error
func IsError(err error) bool {
	_, ok := err.(*Error)
	return ok
}
