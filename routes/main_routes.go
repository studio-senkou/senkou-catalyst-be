package routes

import (
	"senkou-catalyst-be/container"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func InitRoutes(app *fiber.App, deps *container.Container) {
	// Initialize all routes
	// This initiation is use to centralize the route initialization
	// and make it easier to manage dependencies.

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome into the Catalyst API",
		})
	})

	app.Get("/docs/*", swagger.HandlerDefault)

	InitUserRoutes(app, deps.UserController)
	InitAuthRoutes(app, deps.AuthController)
	InitMerchantRoutes(app, deps.MerchantController)
	InitCategoryRoutes(app, deps.CategoryController)
	InitPredefinedCategoryRoutes(app, deps.PredefinedCategoryController)
	InitProductRoutes(app, ProductRouteDependencies{
		ProductController: deps.ProductController,
		UserService:       deps.UserService,
		ProductService:    deps.ProductService,
	})
	InitSubscriptionRoutes(app, deps.SubscriptionController)
	InitPaymentMethodsRoutes(app, deps.PaymentMethodsController)
	InitPaymentRoutes(app, deps.PaymentController)
	InitStorageRoutes(app, deps.StorageController)
}
