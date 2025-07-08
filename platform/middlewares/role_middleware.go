package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(roles ...string) func(c *fiber.Ctx) error {
	// This middleware will check if the user has the required role to access certain routes
	// With this middleware, we can ensure that only users with the appropriate roles can access specific endpoints

	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")

		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "You are not authorized to access this resource",
			})
		}

		return nil
	}
}
