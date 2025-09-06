package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordWebhook struct {
	Content string `json:"content"`
}

type DiscordWebhookBuilder struct {
	WebhookURL string
}

func NewDiscordWebhookBuilder(webhookURL string) *DiscordWebhookBuilder {
	return &DiscordWebhookBuilder{
		WebhookURL: webhookURL,
	}
}

func (b *DiscordWebhookBuilder) Send(payload *DiscordWebhook) error {
	if b.WebhookURL == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(b.WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("webhook returned status code: %d", resp.StatusCode)
	}

	return nil
}

func (b *DiscordWebhookBuilder) SendContent(content string) error {
	payload := &DiscordWebhook{
		Content: content,
	}
	return b.Send(payload)
}

func SendMessage(webhookURL, content string) error {
	builder := NewDiscordWebhookBuilder(webhookURL)
	return builder.SendContent(content)
}
