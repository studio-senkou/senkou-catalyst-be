package config

import (
	"encoding/json"
	"fmt"
	"log"

	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/utils/config"
	"senkou-catalyst-be/utils/response"
	"senkou-catalyst-be/utils/webhook"

	"github.com/gofiber/fiber/v2"
)

func InitFiberConfig() *fiber.Config {
	return &fiber.Config{
		AppName: "Senkou Catalyst API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			resp := response.ErrorResponse{
				Success: false,
				Error: response.ErrorDetail{
					Code:    fiber.StatusInternalServerError,
					Type:    "INTERNAL_ERROR",
					Message: "Internal server error",
				},
			}

			if customErr, ok := err.(*errors.CustomError); ok {
				resp.Error = response.ErrorDetail{
					Code:    customErr.StatusCode(),
					Type:    customErr.ErrorType(),
					Message: customErr.Message,
					Details: customErr.Details,
				}

				if customErr.StatusCode() >= 500 {
					sendLog(c, resp)
				}

				return c.Status(customErr.StatusCode()).JSON(fiber.Map{
					"message": customErr.Message,
					"error":   customErr.Details,
				})
			}

			if fiberErr, ok := err.(*fiber.Error); ok {
				resp.Error.Code = fiberErr.Code
				resp.Error.Message = fiberErr.Message
				resp.Error.Type = "FIBER_ERROR"

				if fiberErr.Code >= 500 {
					sendLog(c, resp)
				}

				return c.Status(fiberErr.Code).JSON(fiber.Map{
					"message": resp.Error.Message,
				})
			}

			resp.Error.Message = err.Error()
			sendLog(c, resp)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
				"error":   resp.Error.Message,
			})
		},
	}
}

func sendLog(c *fiber.Ctx, resp response.ErrorResponse) {
	responseBytes, err := json.Marshal(resp.Error.StackTrace)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		return
	}

	webhookURL := config.GetEnv("WEBHOOK_URL", "")

	builder := webhook.NewDiscordWebhookBuilder(webhookURL)

	message := &webhook.DiscordWebhook{
		Content: fmt.Sprintf(
			"🚨 Error Occured:\n\n**Endpoint**: `%s` \n\n**Error**: %s\n\n**Stack Trace**:\n```%s```",
			c.Path(),
			resp.Error.Details,
			string(responseBytes),
		),
	}

	if err := builder.Send(message); err != nil {
		log.Printf("Failed to send error log to Discord: %v", err)
		return
	}
}
