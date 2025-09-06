package response

import (
	"senkou-catalyst-be/platform/errors"

	"github.com/gofiber/fiber/v2"
)

type ErrorDetail struct {
	Code       int    `json:"code"`
	Type       string `json:"-"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
	StackTrace any    `json:"stack_trace,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"-"`
	Error   ErrorDetail `json:"error"`
	Data    any         `json:"data,omitempty"`
	Meta    any         `json:"meta,omitempty"`
}

func ValidationError(c *fiber.Ctx, message string, fields any) error {
	customErr := errors.Validation(message, fields)

	errors := customErr.Details.(map[string]interface{})["validation_errors"]

	return c.Status(customErr.StatusCode()).JSON(fiber.Map{
		"message": customErr.Message,
		"errors":  errors,
	})
}

func NotFound(c *fiber.Ctx, message string) error {
	customErr := errors.NewCustomError(404, message, "NOT_FOUND", nil)
	return c.Status(customErr.StatusCode()).JSON(customErr)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	customErr := errors.NewCustomError(401, message, "UNAUTHORIZED", nil)
	return c.Status(customErr.StatusCode()).JSON(customErr)
}

func Forbidden(c *fiber.Ctx, message string) error {
	customErr := errors.NewCustomError(403, message, "FORBIDDEN", nil)
	return c.Status(customErr.StatusCode()).JSON(customErr)
}

func BadRequest(c *fiber.Ctx, message string, details any) error {
	customErr := errors.NewCustomError(400, message, "BAD_REQUEST", details)
	return c.Status(customErr.StatusCode()).JSON(customErr)
}

func InternalError(c *fiber.Ctx, message string, details any) error {
	return errors.Internal(message, details)
}

func DatabaseError(c *fiber.Ctx, message string, operation string) error {
	return errors.Database(message, operation)
}

func IsValidationError(err error) bool {
	if customErr, ok := err.(*errors.CustomError); ok {
		return customErr.ErrorType() == "VALIDATION_ERROR"
	}
	return false
}

func IsNotFoundError(err error) bool {
	if customErr, ok := err.(*errors.CustomError); ok {
		return customErr.ErrorType() == "NOT_FOUND"
	}
	return false
}

func IsUnauthorizedError(err error) bool {
	if customErr, ok := err.(*errors.CustomError); ok {
		return customErr.ErrorType() == "UNAUTHORIZED"
	}
	return false
}
