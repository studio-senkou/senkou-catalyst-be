package middlewares

import (
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/config"
	"slices"

	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(roles ...string) func(c *fiber.Ctx) error {
	// This middleware will check if the user has the required role to access certain routes
	// With this middleware, we can ensure that only users with the appropriate roles can access specific endpoints

	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")

		if userID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "You are not authorized to access this resource",
			})
		}

		user := new(models.User)
		err := config.DB.Select("role").Where("id = ?", userID).First(&user).Error

		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "You are not authorized to access this resource",
			})
		}

		if slices.Contains(roles, user.Role) {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You do not have permission to access this resource",
		})
	}
}
