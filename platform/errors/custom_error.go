package errors

import (
	"fmt"
	"net/http"
	"time"
)

type CustomError struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Type      string      `json:"type"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func (e *CustomError) StatusCode() int {
	return e.Code
}

func (e *CustomError) ErrorType() string {
	return e.Type
}

func (e *CustomError) ErrorDetails() interface{} {
	return e.Details
}

func NewCustomError(code int, message, errorType string, details interface{}) *CustomError {
	return &CustomError{
		Code:      code,
		Message:   message,
		Type:      errorType,
		Details:   details,
		Timestamp: time.Now(),
	}
}

func BadRequest(message string, details interface{}) *CustomError {
	return NewCustomError(http.StatusBadRequest, message, "BAD_REQUEST", details)
}

func Unauthorized(message string) *CustomError {
	return NewCustomError(http.StatusUnauthorized, message, "UNAUTHORIZED", nil)
}

func Forbidden(message string) *CustomError {
	return NewCustomError(http.StatusForbidden, message, "FORBIDDEN", nil)
}

func NotFound(message string) *CustomError {
	return NewCustomError(http.StatusNotFound, message, "NOT_FOUND", nil)
}

func Conflict(message string, details interface{}) *CustomError {
	return NewCustomError(http.StatusConflict, message, "CONFLICT", details)
}

func Internal(message string, details interface{}) *CustomError {
	return NewCustomError(http.StatusInternalServerError, message, "INTERNAL_ERROR", details)
}

func Validation(message string, fields interface{}) *CustomError {
	return NewCustomError(http.StatusBadRequest, message, "VALIDATION_ERROR", map[string]interface{}{
		"validation_errors": fields,
	})
}

func Database(message string, operation string) *CustomError {
	return NewCustomError(http.StatusInternalServerError, message, "DATABASE_ERROR", map[string]interface{}{
		"operation": operation,
	})
}

func (e *CustomError) WithDetails(details interface{}) *CustomError {
	e.Details = details
	return e
}

func (e *CustomError) IsClientError() bool {
	return e.Code >= 400 && e.Code < 500
}

func (e *CustomError) IsServerError() bool {
	return e.Code >= 500
}

func (e *CustomError) String() string {
	return fmt.Sprintf("[%s] %s (Code: %d)", e.Type, e.Message, e.Code)
}
