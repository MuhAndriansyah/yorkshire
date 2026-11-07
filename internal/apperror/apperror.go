package apperror

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with an HTTP status and safe message.
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func ErrUnauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message, nil)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message, nil)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, message, nil)
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message, nil)
}

func Internal(err error, message string) *AppError {
	if message == "" {
		message = "Internal server error"
	}
	return New(http.StatusInternalServerError, message, err)
}
