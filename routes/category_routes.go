package routes

import (
	"senkou-catalyst-be/app/controllers"
	"senkou-catalyst-be/platform/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitCategoryRoutes(app *fiber.App, categoryController *controllers.CategoryController) {
	app.Post("/merchants/:merchantID/categories", middlewares.JWTProtected, middlewares.SubscriptionMiddleware("Subscription-Category-Limit"), categoryController.CreateCategory)
	app.Get("/merchants/:merchantID/categories", middlewares.JWTProtected, categoryController.GetCategories)
	app.Put("/merchants/:merchantID/categories/:categoryID", middlewares.JWTProtected, categoryController.UpdateCategory)
	app.Delete("/merchants/:merchantID/categories/:categoryID", middlewares.JWTProtected, categoryController.DeleteCategory)
}
