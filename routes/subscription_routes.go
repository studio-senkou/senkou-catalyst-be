package routes

import (
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/platform/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitSubscriptionRoutes(app *fiber.App, subscriptionController *controllers.SubscriptionController) {
	// Define the routes for subscription
	app.Post("/subscriptions", middlewares.JWTProtected, middlewares.RoleMiddleware("admin"), subscriptionController.CreateSubscription)
	app.Get("/subscriptions", subscriptionController.GetSubscriptions)
	app.Put("/subscriptions/:subID", middlewares.JWTProtected, middlewares.RoleMiddleware("admin"), subscriptionController.UpdateSubscription)
	app.Delete("/subscriptions/:subID", middlewares.JWTProtected, middlewares.RoleMiddleware("admin"), subscriptionController.DeleteSubscription)
	app.Post("/subscriptions/:subID/subscribe", middlewares.JWTProtected, subscriptionController.SubscribeSubscription)
}
