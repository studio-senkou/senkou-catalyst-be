package routes

import (
	"senkou-catalyst-be/app/controllers"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/platform/constants"
	"senkou-catalyst-be/platform/middlewares"

	"github.com/gofiber/fiber/v2"
)

type ProductRouteDependencies struct {
	ProductController *controllers.ProductController
	UserService       services.UserService
	ProductService    services.ProductService
}

func InitProductRoutes(app *fiber.App, deps ProductRouteDependencies) {
	app.Post(
		"/products",
		middlewares.JWTProtected,
		middlewares.SubscriptionMiddleware(constants.SubscriptionProductSlot),
		deps.ProductController.CreateProduct,
	)
	app.Post(
		"/products/:productID/photos",
		middlewares.JWTProtected,
		deps.ProductController.UploadProductPhoto,
	)
	app.Post(
		"/products/:productID/interactions",
		deps.ProductController.SendProductLog,
	)

	app.Get(
		"/products",
		middlewares.JWTProtected,
		middlewares.RoleMiddleware("admin"),
		deps.ProductController.GetAllProducts,
	)
	app.Get(
		"/products/:id",
		deps.ProductController.GetProductByID,
	)
	app.Get(
		"/merchants/:username/popular-products",
		deps.ProductController.PopularProducts,
	)
	app.Get(
		"/merchants/:username/recent-products",
		deps.ProductController.RecentProducts,
	)
	app.Get(
		"/merchants/:username/products",
		deps.ProductController.GetProductByMerchantUsername,
	)

	app.Delete(
		"/products/:productID/photos/*",
		middlewares.JWTProtected,
		deps.ProductController.DeleteProductPhoto,
	)

	route := app.Group(
		"/merchants/:merchantID/products/:productID",
		middlewares.JWTProtected,
		middlewares.OwnershipMiddleware(deps.ProductService, deps.UserService),
	)
	route.Put(
		"/",
		deps.ProductController.UpdateProduct,
	)
	route.Delete(
		"/",
		deps.ProductController.DeleteProduct,
	)
}
