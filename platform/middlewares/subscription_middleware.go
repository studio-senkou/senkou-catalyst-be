package middlewares

import (
	"fmt"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/config"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/utils/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SubscriptionMiddleware(plans ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return response.InternalError(c, "Failed to parse user ID", fmt.Sprintf("Invalid user ID: %v", err.Error()))
		}

		db := config.GetDB()

		userRepo := repositories.NewUserRepository(db)

		user, err := userRepo.FindByID(uint32(userID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve user",
			})
		} else if user != nil && user.Role == "admin" {
			return c.Next()
		}

		subsRepo := repositories.NewSubscriptionRepository(db)

		sub, err := subsRepo.FindActiveSubscriptionByUserID(uint32(userID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve subscription",
			})
		} else if len(sub.Plans) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "User does not have an active subscription",
			})
		}

		for _, plan := range sub.Plans {
			if contains(plans, plan.Name) {
				if hasAccess, err := hasAccess(&plan, user); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to verify subscription plan access",
					})
				} else if hasAccess {
					return c.Next()
				} else {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "User does not have access to this feature",
					})
				}
			}
		}

		return c.Next()
	}
}

func hasAccess(plan *models.SubscriptionPlan, user *models.User) (bool, error) {
	db := config.GetDB()

	productRepository := repositories.NewProductRepository(db)
	categoryRepository := repositories.NewCategoryRepository(db)

	switch plan.Name {
	case "Subscription-Product-Slot":

		if products, err := productRepository.FindProductsByMerchantID(user.Merchants[0].ID); err != nil {
			return false, err
		} else {
			max, err := strconv.Atoi(plan.Value)
			if err != nil {
				return false, err
			}
			if len(products) >= max {
				return false, nil
			}
		}

		return true, nil

	case "Subscription-Category-Limit":

		if categories, err := categoryRepository.FindAllCategoriesByMerchantID(user.Merchants[0].ID); err != nil {
			return false, err
		} else {
			max, err := strconv.Atoi(plan.Value)
			if err != nil {
				return false, err
			}
			if len(categories) >= max {
				return false, nil
			}
		}

		return true, nil
	}

	return false, nil
}

func contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
