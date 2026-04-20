package model

import "net/http"

// AppError is a business-level error that carries an HTTP status code.
// Usecases return AppError; the global ErrorHandler maps it to JSON response.
type AppError struct {
	Code    int    // HTTP status code
	Field   string // optional — for field-level errors
	Message string
}

func (e *AppError) Error() string { return e.Message }

func ErrConflict(field, message string) *AppError {
	return &AppError{Code: http.StatusConflict, Field: field, Message: message}
}

func ErrBadRequest(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

func ErrUnauthorized(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message}
}

func ErrForbidden(message string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: message}
}

func ErrNotFound(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

func ErrInternal(message string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}
