package config

import (
	"senkou-catalyst-be/utils/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitCorsConfig() fiber.Handler {
	origins := config.GetEnv("APP_ALLOWED_ORIGINS", "http://localhost:5173")

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Range",
		MaxAge:           86400,
	})
}
