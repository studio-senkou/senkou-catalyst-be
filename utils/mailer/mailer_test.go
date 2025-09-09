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

func TestParseTemplate_Valid(t *testing.T) {
	// Create a temp template file
	tmplPath := "test-template.html"
	tmplContent := "Hello, {{.Name}}!"
	os.WriteFile("utils/mailer/templates/"+tmplPath, []byte(tmplContent), 0644)
	defer os.Remove("utils/mailer/templates/" + tmplPath)

	result, err := parseTemplate(tmplPath, map[string]string{"Name": "World"})
	if err != nil {
		t.Errorf("parseTemplate failed: %v", err)
	}
	if result != "Hello, World!" {
		t.Errorf("parseTemplate output wrong: %s", result)
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
