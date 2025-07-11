package errors

import "net/http"

type ValidationError struct {
	*BaseError
	Fields map[string]string `json:"fields"`
}

func NewValidationError(message string, fields map[string]string) *ValidationError {
	return &ValidationError{
		BaseError: &BaseError{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: message,
			ErrorType:    "VALIDATION_ERROR",
		},
		Fields: fields,
	}
}
