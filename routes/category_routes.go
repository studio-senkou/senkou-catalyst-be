package routes

import (
	"senkou-catalyst-be/app/controllers"
	"senkou-catalyst-be/platform/constants"
	"senkou-catalyst-be/platform/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitCategoryRoutes(app *fiber.App, categoryController *controllers.CategoryController) {
	app.Post(
		"/merchants/:merchantID/categories",
		middlewares.JWTProtected,
		middlewares.SubscriptionMiddleware(constants.SubscriptionCategoryLimit),
		categoryController.CreateCategory,
	)
	app.Get(
		"/merchants/:merchantID/categories",
		middlewares.JWTProtected,
		categoryController.GetCategories,
	)
	app.Put(
		"/merchants/:merchantID/categories/:categoryID",
		middlewares.JWTProtected,
		categoryController.UpdateCategory,
	)
	app.Delete(
		"/merchants/:merchantID/categories/:categoryID",
		middlewares.JWTProtected,
		categoryController.DeleteCategory,
	)

	// Category management using merchant username
	// This will allow user to manage categories without needing merchant ID
	// This is useful for client where merchant ID is not known or just knowing merchant username
	app.Post(
		"/merchants/username/:username/categories",
		middlewares.JWTProtected,
		middlewares.SubscriptionMiddleware(constants.SubscriptionCategoryLimit),
		categoryController.CreateCategoryWithMerchantUsername,
	)

	app.Get(
		"/merchants/username/:username/categories",
		categoryController.GetCategoriesByMerchantUsername,
	)

}
