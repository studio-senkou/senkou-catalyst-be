package mailer

import (
	"os"
	"testing"
)

func TestNewMail(t *testing.T) {
	mail := NewMail("to@example.com", "Subject", "Body")
	if mail.To != "to@example.com" || mail.Subject != "Subject" || mail.Body != "Body" || mail.IsHTML {
		t.Errorf("NewMail did not set fields correctly: %+v", mail)
	}
}

func TestNewHTMLMail(t *testing.T) {
	mail := NewHTMLMail("to@example.com", "Subject", "<b>HTML</b>")
	if mail.To != "to@example.com" || mail.Subject != "Subject" || mail.HTMLBody != "<b>HTML</b>" || !mail.IsHTML {
		t.Errorf("NewHTMLMail did not set fields correctly: %+v", mail)
	}
}

func TestNewMailFromTemplate_InvalidPath(t *testing.T) {
	_, err := NewMailFromTemplate("to@example.com", "Subject", "not-exist.html", nil)
	if err == nil {
		t.Error("Expected error for invalid template path, got nil")
	}
}

func TestNewMailFromTemplate_ValidTemplate(t *testing.T) {
	data := map[string]interface{}{
		"ActivationLink": "https://example.com/activate",
		"SupportEmail":   "support@example.com",
	}

	mail, err := NewMailFromTemplate("to@example.com", "Subject", "account-activation.html", data)
	if err != nil {
		t.Errorf("NewMailFromTemplate failed: %v", err)
	}

	if mail.To != "to@example.com" {
		t.Errorf("Expected To: to@example.com, got: %s", mail.To)
	}

	if !mail.IsHTML {
		t.Error("Mail should be HTML")
	}

	if mail.HTMLBody == "" {
		t.Error("HTMLBody should not be empty")
	}
}

func TestEmbeddedTemplateManagerParsing(t *testing.T) {
	tm := NewTemplateManager()

	data := map[string]interface{}{
		"ActivationLink": "https://example.com/activate",
		"SupportEmail":   "support@example.com",
	}

	result, err := tm.ParseTemplate("account-activation.html", data)
	if err != nil {
		t.Errorf("ParseTemplate failed: %v", err)
	}
	if result == "" {
		t.Error("ParseTemplate returned empty result")
	}
}

func TestMail_Send_InvalidConfig(t *testing.T) {
	mail := NewMail("to@example.com", "Subject", "Body")
	os.Setenv("MAILER_HOST", "")
	os.Setenv("MAILER_PORT", "")
	os.Setenv("MAILER_USERNAME", "")
	os.Setenv("MAILER_PASSWORD", "")
	err := mail.Send()
	if err == nil {
		t.Error("Expected error for missing SMTP config, got nil")
	}
}

func TestMail_Send_EmptyFields(t *testing.T) {
	mail := &Mail{To: "", Subject: "", Body: ""}
	os.Setenv("MAILER_HOST", "smtp.example.com")
	os.Setenv("MAILER_PORT", "587")
	os.Setenv("MAILER_USERNAME", "user")
	os.Setenv("MAILER_PASSWORD", "pass")
	err := mail.Send()
	if err == nil {
		t.Error("Expected error for empty to/subject, got nil")
	}
}
