package errors

import "net/http"

type InternalError struct {
	*BaseError
}

func NewInternalError(message string, details any) *InternalError {
	return &InternalError{
		BaseError: &BaseError{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: message,
			ErrorType:    "INTERNAL_ERROR",
			ErrorDetails: details,
		},
	}
}

type DatabaseError struct {
	*BaseError
}

func NewDatabaseError(message string, operation string) *DatabaseError {
	return &DatabaseError{
		BaseError: &BaseError{
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: message,
			ErrorType:    "DATABASE_ERROR",
			ErrorDetails: map[string]interface{}{
				"operation": operation,
			},
		},
	}
}
