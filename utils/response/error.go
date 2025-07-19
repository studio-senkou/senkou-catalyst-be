package response

import (
	"senkou-catalyst-be/platform/errors"

	"github.com/gofiber/fiber/v2"
)

func ValidationError(c *fiber.Ctx, message string, fields any) error {
	return errors.NewValidationError(message, fields)
}

func NotFound(c *fiber.Ctx, resource string) error {
	return errors.NewNotFoundError(resource)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return errors.NewUnauthorizedError(message)
}

func Forbidden(c *fiber.Ctx, message string) error {
	return errors.NewForbiddenError(message)
}

func BadRequest(c *fiber.Ctx, message string, details any) error {
	return errors.NewBadRequestError(message, details)
}

func InternalError(c *fiber.Ctx, message string, details any) error {
	return errors.NewInternalError(message, details)
}

func DatabaseError(c *fiber.Ctx, message string, operation string) error {
	return errors.NewDatabaseError(message, operation)
}

func IsValidationError(err error) bool {
	_, ok := err.(*errors.ValidationError)
	return ok
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*errors.NotFoundError)
	return ok
}

func IsUnauthorizedError(err error) bool {
	_, ok := err.(*errors.UnauthorizedError)
	return ok
}
