package middlewares

import (
	"fmt"
	"senkou-catalyst-be/utils"

	"github.com/gofiber/fiber/v2"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(utils.GetEnv("AUTH_SECRET", ""))

	if len(jwtSecret) == 0 {
		panic("JWT secret is not set. Please set the AUTH_SECRET environment variable.")
	}
}

var JWTProtected fiber.Handler = func(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to access this resource",
		})
	}

	var authToken string
	
	_, err := fmt.Sscanf(token, "Bearer %s", &authToken)

	if err != nil || authToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid authorization header format",
		})
	}

	claims, err := utils.ValidateToken(authToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
			"error":   err.Error(),
		})
	}

	c.Locals("user_id", claims["payload"])

	return c.Next()
}
