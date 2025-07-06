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

// @title Fiber API with Service-Repository
// @version 1.0
// @description Fiber REST API using service-repo pattern
// @host localhost:3000
// @BasePath /

// @schemes http
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

	err := app.Listen(fmt.Sprintf(":%s", utils.GetEnv("APP_PORT", "8080")))

	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
