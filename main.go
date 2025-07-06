package main

import (
	"fmt"
	"senkou-catalyst-be/config"
	"senkou-catalyst-be/routes"
	"senkou-catalyst-be/utils"

	_ "senkou-catalyst-be/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @Title Catalyst API Swagger
// @Version 1.0
// @Description This is the API documentation for the Catalyst application.
// @Contact.name Catalyst Team
// @Host localhost:8080
// @BasePath /
// @Schemes http
func main() {
	appConfig := config.InitFiberConfig()
	app := fiber.New(*appConfig)

	config.ConnectDB()

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome into the Catalyst API",
		})
	})

	routes.UserRoutes(app)
	routes.AuthRoutes(app)
	routes.MerchantRoutes(app)

	err := app.Listen(fmt.Sprintf(":%s", utils.GetEnv("APP_PORT", "8080")))

	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
