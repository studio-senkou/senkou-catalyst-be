package routes

import (
	"senkou-catalyst-be/app/controllers"

	goth "senkou-catalyst-be/integrations/goth"

	"github.com/gofiber/fiber/v2"
)

func InitOAuthRoutes(app *fiber.App, oauthController *controllers.OAuthController) {
	app.Get("/auth/:provider", goth.BeginAuthHandler)
	app.Get("/auth/:provider/callback", oauthController.GoogleCallback)
}
