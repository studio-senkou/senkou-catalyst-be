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

// This middleware checks if the request has a valid JWT token that provided by the user
// If the token is valid, it extracts the payload that contain the user ID and stores it in the context
// If the token is invalid or missing, it returns a Uauthorized response
// Generally used to protect routes that require authentication
var JWTProtected fiber.Handler = func(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	// Check if the user provide with the Authorization header
	// If the header is not present, we return an Uauthorized response
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to access this resource",
		})
	}

	var authToken string

	_, err := fmt.Sscanf(token, "Bearer %s", &authToken)

	// If the token is not in the correct format, we return an Uauthorized response
	// The correct format is "Bearer <token>"
	if err != nil || authToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid authorization header format",
		})
	}

	// Validate the token and claims or extract the payload inside the token
	// If the token is valid, we extract the payload that contain the user ID and store
	// it in the context for further processing
	claims, err := utils.ValidateToken(authToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Cannot continue to process request due to invalid token",
			"error":   fmt.Sprintf("Token validation failed: %s", err.Error()),
		})
	}

	c.Locals("userID", claims["payload"])

	return c.Next()
}
