package mailer

import (
	"testing"
)

func TestTemplateManager(t *testing.T) {
	tm := NewTemplateManager()

	if !tm.TemplateExists("account-activation.html") {
		t.Error("account-activation.html template should exist")
	}

	if tm.TemplateExists("non-existent.html") {
		t.Error("non-existent.html template should not exist")
	}

	templates := tm.GetAvailableTemplates()
	if len(templates) == 0 {
		t.Error("should have at least one template")
	}

	found := false
	for _, template := range templates {
		if template == "account-activation.html" {
			found = true
			break
		}
	}
	if !found {
		t.Error("account-activation.html should be in available templates")
	}
}

func TestTemplateManagerParse(t *testing.T) {
	tm := NewTemplateManager()

	data := map[string]interface{}{
		"ActivationLink": "https://example.com/activate",
		"SupportEmail":   "support@example.com",
	}

	result, err := tm.ParseTemplate("account-activation.html", data)
	if err != nil {
		t.Errorf("Failed to parse template: %v", err)
	}

	if result == "" {
		t.Error("Parsed template should not be empty")
	}

	_, err = tm.ParseTemplate("non-existent.html", data)
	if err == nil {
		t.Error("Should return error for non-existent template")
	}
}

func TestMailBuilder(t *testing.T) {
	builder := NewMailBuilder()

	mail, err := builder.
		To("test@example.com").
		Subject("Test Subject").
		PlainBody("Test body").
		Build()

	if err != nil {
		t.Errorf("Failed to build mail: %v", err)
	}

	if mail.To != "test@example.com" {
		t.Errorf("Expected To: test@example.com, got: %s", mail.To)
	}

	if mail.Subject != "Test Subject" {
		t.Errorf("Expected Subject: Test Subject, got: %s", mail.Subject)
	}

	if mail.Body != "Test body" {
		t.Errorf("Expected Body: Test body, got: %s", mail.Body)
	}

	_, err = NewMailBuilder().Build()
	if err == nil {
		t.Error("Should return validation error for empty mail")
	}
}
