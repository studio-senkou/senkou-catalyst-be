package config

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/utils/config"
	"senkou-catalyst-be/utils/response"
	"senkou-catalyst-be/utils/webhook"

	"github.com/gofiber/fiber/v2"
)

func InitFiberConfig() *fiber.Config {
	return &fiber.Config{
		AppName:                 "Senkou Catalyst API",
		DisableStartupMessage:   true,
		CaseSensitive:           true,
		StrictRouting:           true,
		EnableTrustedProxyCheck: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			log.Println("Error Handler Invoked")
			log.Printf("Error: %v", err)

			stackTrace := ""
			if customErr, ok := err.(*errors.CustomError); ok && customErr.StatusCode() >= 500 {
				stackTrace = string(debug.Stack())
			} else if fiberErr, ok := err.(*fiber.Error); ok && fiberErr.Code >= 500 {
				stackTrace = string(debug.Stack())
			} else if _, ok := err.(*errors.CustomError); !ok {
				if _, ok := err.(*fiber.Error); !ok {
					stackTrace = string(debug.Stack())
				}
			}

			resp := response.ErrorResponse{
				Success: false,
				Error: response.ErrorDetail{
					Code:       fiber.StatusInternalServerError,
					Type:       "INTERNAL_ERROR",
					Message:    "Internal server error",
					StackTrace: stackTrace,
				},
			}

			if customErr, ok := err.(*errors.CustomError); ok {
				resp.Error = response.ErrorDetail{
					Code:       customErr.StatusCode(),
					Type:       customErr.ErrorType(),
					Message:    customErr.Message,
					Details:    customErr.Details,
					StackTrace: stackTrace,
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
				resp.Error.StackTrace = stackTrace

				if fiberErr.Code >= 500 {
					sendLog(c, resp)
				}

				return c.Status(fiberErr.Code).JSON(fiber.Map{
					"message": resp.Error.Message,
				})
			}

			resp.Error.Message = err.Error()
			resp.Error.StackTrace = stackTrace
			sendLog(c, resp)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error brok",
				"error":   resp.Error.Message,
			})
		},
	}
}

func sendLog(c *fiber.Ctx, resp response.ErrorResponse) {
	stackTraceStr := ""
	if resp.Error.StackTrace != nil {
		if stackBytes, ok := resp.Error.StackTrace.(string); ok {
			lines := strings.Split(stackBytes, "\n")
			if len(lines) > 20 {
				stackTraceStr = strings.Join(lines[:20], "\n") + "\n... (truncated)"
			} else {
				stackTraceStr = stackBytes
			}
		}
	}

	webhookURL := config.GetEnv("WEBHOOK_URL", "")
	webhookEnabled := config.GetEnv("WEBHOOK_ENABLED", "false")

	if webhookURL == "" || webhookEnabled != "true" {
		log.Printf("🚨 Error Occured:\nEndpoint: %s %s\nError: %s\nDetails: %v\nStack Trace:\n%s",
			c.Method(), c.Path(), resp.Error.Message, resp.Error.Details, stackTraceStr)
		return
	}

	builder := webhook.NewDiscordWebhookBuilder(webhookURL)

	content := fmt.Sprintf(
		"🚨 **Server Error**\n\n**Endpoint**: `%s %s`\n**Error**: %s\n**Details**: %v",
		c.Method(),
		c.Path(),
		resp.Error.Message,
		resp.Error.Details,
	)

	if len(content) < 1500 && stackTraceStr != "" {
		content += fmt.Sprintf("\n\n**Stack Trace**:\n```go\n%s\n```", stackTraceStr)
	} else if stackTraceStr != "" {
		content += "\n\n⚠️ *Stack trace omitted (message too long)*"
	}

	message := &webhook.DiscordWebhook{
		Content: content,
	}

	if err := builder.Send(message); err != nil {
		log.Printf("Failed to send error log to Discord: %v", err)
		log.Printf("🚨 Error Occured:\nEndpoint: %s %s\nError: %s\nDetails: %v\nStack Trace:\n%s",
			c.Method(), c.Path(), resp.Error.Message, resp.Error.Details, stackTraceStr)
		return
	}
}
