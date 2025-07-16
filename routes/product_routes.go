package routes

import (
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/platform/middlewares"
	"senkou-catalyst-be/services"

	"github.com/gofiber/fiber/v2"
)

type ProductRouteDependencies struct {
	ProductController *controllers.ProductController
	UserService       services.UserService
	ProductService    services.ProductService
}

func InitProductRoutes(app *fiber.App, deps ProductRouteDependencies) {
	app.Post("/products", middlewares.JWTProtected, deps.ProductController.CreateProduct)
	app.Get("/products", middlewares.JWTProtected, middlewares.RoleMiddleware("admin"), deps.ProductController.GetAllProducts)
	app.Get("/products/:id", deps.ProductController.GetProductByID)
	app.Get("/merchants/:merchantID/products", deps.ProductController.GetProductByMerchant)

	route := app.Group(
		"/merchants/:merchantID/products/:productID",
		middlewares.JWTProtected,
		middlewares.OwnershipMiddleware(deps.ProductService, deps.UserService),
	)

	route.Put("/", deps.ProductController.UpdateProduct)
	route.Delete("/", deps.ProductController.DeleteProduct)
}
