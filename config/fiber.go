package config

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"senkou-catalyst-be/errors"
	"senkou-catalyst-be/utils"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Success bool           `json:"-"`
	Error   ErrorDetail    `json:"error"`
	Data    any            `json:"data,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
}

type ErrorDetail struct {
	Code    int            `json:"-"`
	Type    string         `json:"-"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func SendLog(errMsg string) {
	discordWebhookURL := utils.GetEnv("DISCORD_WEBHOOK_URL", "")

	if discordWebhookURL == "" {
		return
	}

	var errorResponse ErrorResponse
	json.Unmarshal([]byte(errMsg), &errorResponse)

	embed := map[string]any{
		"title":       "🚨 Error in Senkou Catalyst API",
		"color":       getColorByErrorType(errorResponse.Error.Type),
		"description": "An error occurred in the application",
		"fields": []map[string]any{
			{
				"name":   "Error Type",
				"value":  errorResponse.Error.Type,
				"inline": true,
			},
			{
				"name":   "Status Code",
				"value":  errorResponse.Error.Code,
				"inline": true,
			},
			{
				"name":   "Message",
				"value":  errorResponse.Error.Message,
				"inline": false,
			},
		},
		"timestamp": time.Now().Format(time.RFC3339),
		"footer": map[string]any{
			"text":     "Senkou Catalyst API",
			"icon_url": "https://via.placeholder.com/16x16/ff0000/ffffff?text=!",
		},
	}

	if len(errorResponse.Error.Details) > 0 {
		detailsJson, _ := json.MarshalIndent(errorResponse.Error.Details, "", "  ")
		embed["fields"] = append(embed["fields"].([]map[string]any), map[string]any{
			"name":   "Details",
			"value":  "```json\n" + string(detailsJson) + "\n```",
			"inline": false,
		})
	}

	payload := map[string]any{
		"embeds": []map[string]any{embed},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal webhook payload: %v", err)
		return
	}

	go func() {
		response, postErr := http.Post(discordWebhookURL, "application/json", bytes.NewBuffer(payloadBytes))
		if postErr != nil {
			log.Printf("Failed to send webhook log: %v", postErr)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
			log.Printf("Discord webhook returned status: %d", response.StatusCode)
		}
	}()
}

func InitFiberConfig() *fiber.Config {
	return &fiber.Config{
		AppName: "Senkou Catalyst API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Error: %v", err)

			response := ErrorResponse{
				Success: false,
				Error: ErrorDetail{
					Code:    fiber.StatusInternalServerError,
					Type:    "INTERNAL_ERROR",
					Message: "Internal server error",
				},
			}

			if appErr, ok := err.(errors.AppError); ok {
				response.Error = ErrorDetail{
					Code:    appErr.Code(),
					Type:    appErr.Type(),
					Message: appErr.Error(),
					Details: appErr.Details(),
				}

				if validationErr, ok := err.(*errors.ValidationError); ok {
					response.Error.Details = map[string]any{
						"fields": validationErr.Fields,
					}
				}

				if appErr.Code() >= 500 {
					sendErrorLog(response)
				}
				return c.Status(appErr.Code()).JSON(response)
			}

			if fiberErr, ok := err.(*fiber.Error); ok {
				response.Error.Code = fiberErr.Code
				response.Error.Message = fiberErr.Message
				response.Error.Type = "FIBER_ERROR"

				if fiberErr.Code >= 500 {
					sendErrorLog(response)
				}

				return c.Status(fiberErr.Code).JSON(response)
			}

			response.Error.Message = err.Error()

			sendErrorLog(response)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		},
	}
}

func sendErrorLog(response ErrorResponse) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		return
	}
	SendLog(string(responseBytes))
}

func getColorByErrorType(errorType string) int {
	switch errorType {
	case "VALIDATION_ERROR":
		return 16776960
	case "NOT_FOUND_ERROR":
		return 16753920
	case "UNAUTHORIZED_ERROR", "FORBIDDEN_ERROR":
		return 16711680
	case "DATABASE_ERROR":
		return 8388736
	case "BUSINESS_ERROR":
		return 16744448
	default:
		return 15158332
	}
}
