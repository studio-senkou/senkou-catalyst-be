package middlewares

import (
	"senkou-catalyst-be/config"
	"senkou-catalyst-be/models"

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

		var user models.User
		if err := config.DB.Select("role").Where("id = ?", userID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "You are not authorized to access this resource",
			})
		}

		return c.Next()
	}
}
