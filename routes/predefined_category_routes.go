package routes

import (
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/platform/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitPredefinedCategoryRoutes(app *fiber.App, predefinedCategoryController *controllers.PredefinedCategoryController) {
	app.Get("/predefined-categories", predefinedCategoryController.GetPredefinedCategories)

	pdRoute := app.Group("/predefined-categories", middlewares.JWTProtected, middlewares.RoleMiddleware("admin"))

	pdRoute.Post("/", predefinedCategoryController.StoreCategory)
	pdRoute.Put("/:pcID", predefinedCategoryController.UpdatePredefinedCategory)
	pdRoute.Delete("/:pcID", predefinedCategoryController.DeletePredefinedCategory)
}
