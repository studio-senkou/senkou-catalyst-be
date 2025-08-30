package main

import (
	"fmt"
	"log"
	"senkou-catalyst-be/container"
	"senkou-catalyst-be/platform/config"
	"senkou-catalyst-be/routes"

	environment "senkou-catalyst-be/utils/config"

	_ "senkou-catalyst-be/docs"

	"github.com/gofiber/fiber/v2"
)

// @Title Catalyst API Documentation
// @Description This is the API documentation for the Catalyst application.
// @Version 1.0
// @Contact.name Catalyst Team
// @Host localhost:8080
// @BasePath /
// @Schemes http
func main() {
	appConfig := config.InitFiberConfig()
	app := fiber.New(*appConfig)

	// Middlewares
	app.Use(config.InitCorsConfig())
	app.Use(config.InitLogger())

	config.ConnectDB()

	deps, err := container.InitializeContainer()

	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	routes.InitRoutes(app, deps)

	err = app.Listen(fmt.Sprintf(":%s", environment.GetEnv("APP_PORT", "8080")))

	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
