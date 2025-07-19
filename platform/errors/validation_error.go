package errors

import "net/http"

type ValidationError struct {
	*BaseError
	Fields any `json:"fields"`
}

func NewValidationError(message string, fields any) *ValidationError {
	return &ValidationError{
		BaseError: &BaseError{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: message,
			ErrorType:    "VALIDATION_ERROR",
		},
		Fields: fields,
	}
}
