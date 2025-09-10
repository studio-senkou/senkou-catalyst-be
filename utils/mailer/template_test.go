package mailer

import (
	"testing"
)

func TestNewTemplateManager(t *testing.T) {
	tm := NewTemplateManager()

	if tm == nil {
		t.Fatal("NewTemplateManager returned nil")
	}

	if tm.templates == nil {
		t.Fatal("templates map is nil")
	}

	if len(tm.templates) == 0 {
		t.Fatal("templates map is empty")
	}

	if _, exists := tm.templates["account-activation.html"]; !exists {
		t.Error("account-activation.html template not found")
	}
}

func TestTemplateManager_ParseTemplate(t *testing.T) {
	tm := NewTemplateManager()

	tests := []struct {
		name         string
		templateName string
		data         interface{}
		expectError  bool
		expectEmpty  bool
	}{
		{
			name:         "valid template with data",
			templateName: "account-activation.html",
			data:         map[string]string{"name": "John"},
			expectError:  false,
			expectEmpty:  false,
		},
		{
			name:         "valid template with nil data",
			templateName: "account-activation.html",
			data:         nil,
			expectError:  false,
			expectEmpty:  false,
		},
		{
			name:         "non-existent template",
			templateName: "non-existent.html",
			data:         nil,
			expectError:  false,
			expectEmpty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tm.ParseTemplate(tt.templateName, tt.data)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectEmpty && result != "" {
				t.Error("expected empty result but got content")
			}

			if !tt.expectEmpty && !tt.expectError && result == "" {
				t.Error("expected non-empty result but got empty string")
			}
		})
	}
}

func TestTemplateManager_GetAvailableTemplates(t *testing.T) {
	tm := NewTemplateManager()

	templates := tm.GetAvailableTemplates()

	if len(templates) == 0 {
		t.Fatal("no templates returned")
	}

	found := false
	for _, template := range templates {
		if template == "account-activation.html" {
			found = true
			break
		}
	}

	if !found {
		t.Error("account-activation.html not found in available templates")
	}
}

func TestTemplateManager_ParseTemplate_InvalidTemplate(t *testing.T) {
	tm := NewTemplateManager()

	tm.templates["invalid"] = "{{.InvalidSyntax"

	_, err := tm.ParseTemplate("invalid", nil)
	if err == nil {
		t.Error("expected parse error for invalid template syntax")
	}
}

func TestTemplateManager_ParseTemplate_ExecuteError(t *testing.T) {
	tm := NewTemplateManager()

	tm.templates["execute-error"] = "{{.NonExistentField.SubField}}"

	_, err := tm.ParseTemplate("execute-error", map[string]string{})
	if err == nil {
		t.Error("expected execution error for template with invalid field access")
	}
}
