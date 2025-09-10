package mailer

import (
	"fmt"
)

type MailerService struct {
	templateManager *TemplateManager
	smtpConfig      *SMTPConfig
}

func NewMailerService() (*MailerService, error) {
	config, err := NewSMTPConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SMTP config: %w", err)
	}

	return &MailerService{
		templateManager: NewTemplateManager(),
		smtpConfig:      config,
	}, nil
}

func (ms *MailerService) NewBuilder() *MailBuilder {
	return NewMailBuilder()
}

func (ms *MailerService) SendTemplate(to, subject, templateName string, data interface{}) error {
	mail, err := ms.NewBuilder().
		To(to).
		Subject(subject).
		Template(templateName, data).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build mail: %w", err)
	}

	return mail.SendWithConfig(ms.smtpConfig)
}

func (ms *MailerService) SendPlainText(to, subject, body string) error {
	mail, err := ms.NewBuilder().
		To(to).
		Subject(subject).
		PlainBody(body).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build mail: %w", err)
	}

	return mail.SendWithConfig(ms.smtpConfig)
}

func (ms *MailerService) SendHTML(to, subject, htmlBody string) error {
	mail, err := ms.NewBuilder().
		To(to).
		Subject(subject).
		HTMLBody(htmlBody).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build mail: %w", err)
	}

	return mail.SendWithConfig(ms.smtpConfig)
}

func (ms *MailerService) GetAvailableTemplates() []string {
	return ms.templateManager.GetAvailableTemplates()
}

func (ms *MailerService) TemplateExists(templateName string) bool {
	return ms.templateManager.TemplateExists(templateName)
}
