package routes

import (
	"senkou-catalyst-be/app/controllers"
	"senkou-catalyst-be/platform/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitPredefinedCategoryRoutes(app *fiber.App, PDController *controllers.PredefinedCategoryController) {
	app.Get(
		"/predefined-categories",
		PDController.GetPredefinedCategories,
	)

	PDRoute := app.Group(
		"/predefined-categories",
		middlewares.JWTProtected,
		middlewares.RoleMiddleware("admin"),
	)

	PDRoute.Post(
		"/",
		PDController.StoreCategory,
	)
	PDRoute.Put(
		"/:pcID",
		PDController.UpdatePredefinedCategory,
	)
	PDRoute.Delete("/:pcID",
		PDController.DeletePredefinedCategory,
	)
}
