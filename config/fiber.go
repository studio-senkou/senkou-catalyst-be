package config

import "github.com/gofiber/fiber/v2"

func InitFiberConfig() *fiber.Config {
	return &fiber.Config{
		AppName: "Senkou Catalyst API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
				"error":   err.Error(),
			})
		},
	}
}
