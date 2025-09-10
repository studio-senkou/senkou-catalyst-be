package mailer

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed templates/account-activation.html
var accountActivationTemplate string

type TemplateManager struct {
	templates map[string]string
}

func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		templates: map[string]string{
			"account-activation.html": accountActivationTemplate,
			// Add more templates here as needed
			// "password-reset.html": passwordResetTemplate,
			// "welcome.html": welcomeTemplate,
		},
	}
}

func (tm *TemplateManager) ParseTemplate(templateName string, data interface{}) (string, error) {
	templateContent, exists := tm.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	tmpl, err := template.New(templateName).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

func (tm *TemplateManager) GetAvailableTemplates() []string {
	templates := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		templates = append(templates, name)
	}
	return templates
}

func (tm *TemplateManager) TemplateExists(templateName string) bool {
	_, exists := tm.templates[templateName]
	return exists
}
