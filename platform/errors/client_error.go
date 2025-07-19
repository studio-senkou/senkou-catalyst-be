package errors

import (
	"net/http"
)

type BadRequestError struct {
	*BaseError
}

func NewBadRequestError(message string, details any) *BadRequestError {
	return &BadRequestError{
		BaseError: &BaseError{
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: message,
			ErrorType:    "BUSINESS_ERROR",
			ErrorDetails: details,
		},
	}
}

type UnauthorizedError struct {
	*BaseError
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		BaseError: &BaseError{
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: message,
			ErrorType:    "UNAUTHORIZED_ERROR",
		},
	}
}

type ForbiddenError struct {
	*BaseError
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{
		BaseError: &BaseError{
			ErrorCode:    http.StatusForbidden,
			ErrorMessage: message,
			ErrorType:    "FORBIDDEN_ERROR",
		},
	}
}

type NotFoundError struct {
	*BaseError
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		BaseError: &BaseError{
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: message,
			ErrorType:    "NOT_FOUND_ERROR",
		},
	}
}
