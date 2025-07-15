package middlewares

import (
	"fmt"
	"senkou-catalyst-be/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func OwnershipMiddleware(productService services.ProductService, userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
		userID, err := strconv.ParseUint(userIDStr, 10, 64)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to parse user ID",
				"error":   err.Error(),
			})
		}

		productID := c.Params("productID")

		if productID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Product ID is required",
			})
		}

		if err := productService.VerifyProductOwnership(productID, uint32(userID)); err != nil {

			// Validate if the user is an Administrator
			if isAdmin, err := userService.VerifyIsAnAdministrator(uint32(userID)); err != nil || !isAdmin {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "You do not have permission to access this resource",
				})
			}

			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not own this product",
			})
		}

		return c.Next()
	}
}
