package webhook

import (
	"os"
	"testing"
)

func TestDiscordWebhook(t *testing.T) {

	t.Run("Should not be able to send message with invalid webhook URL", func(t *testing.T) {
		webhookURL := ""
		content := "Test message from unit test"

		err := SendMessage(webhookURL, content)
		if err != nil {
			t.Logf("Expected error occurred: %v", err)
		}
	})

	t.Run("Should be able to create discord webhook builder", func(t *testing.T) {
		builder := NewDiscordWebhookBuilder("https://discord.com/api/webhooks/mywebhookid/mywebhooktoken")
		if builder == nil {
			t.Error("Failed to create Discord webhook builder")
		}

		t.Logf("Created Discord webhook builder: %+v", builder)
	})

	t.Run("Should be able to send message with valid webhook URL", func(t *testing.T) {
		webhookURL := os.Getenv("WEBHOOK_URL")
		content := "Test message from unit test"

		if webhookURL == "" {
			t.Skip("WEBHOOK_URL is not set, skipping test")
		}

		err := SendMessage(webhookURL, content)
		if err != nil {
			t.Errorf("Failed to send message: %v", err)
		} else {
			t.Log("Message sent successfully")
		}
	})

	t.Run("Should be able to send message using builder", func(t *testing.T) {
		webhookURL := os.Getenv("WEBHOOK_URL")
		content := "Test message from builder in unit test"

		if webhookURL == "" {
			t.Skip("WEBHOOK_URL is not set, skipping test")
		}

		builder := NewDiscordWebhookBuilder(webhookURL)
		err := builder.SendContent(content)
		if err != nil {
			t.Errorf("Failed to send message using builder: %v", err)
		}
	})

}
