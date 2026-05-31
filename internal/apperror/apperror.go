// package apperror contains an application error for use in the transport layer
package apperror

import "fmt"

// AppError is an error for use at the transport layer (http handlers). It integrates with the
// server handlers so that any error returned from a handler of this type will be json marshalled.
type AppError struct {
	Msg  any `json:"error"`
	Code int `json:"-"`
}

// New return a new AppError instance
func New(message any, code int) *AppError {
	return &AppError{
		Msg:  message,
		Code: code,
	}
}

// New return a new AppError instance from an error
func NewFromError(err error, code int) *AppError {
	return New(err.Error(), code)
}

func (e *AppError) Error() string {
	return fmt.Sprintf("app error with code %d: %v", e.Code, e.Msg)
}
