package mailer

import (
	"crypto/tls"
	"errors"
	"fmt"

	"gopkg.in/gomail.v2"
)

type Mail struct {
	To         string
	Subject    string
	Body       string
	HTMLBody   string
	IsHTML     bool
	buildError error
}

func NewMail(to, subject, body string) *Mail {
	return &Mail{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	}
}

func NewHTMLMail(to, subject, htmlBody string) *Mail {
	return &Mail{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlBody,
		IsHTML:   true,
	}
}

func NewMailFromTemplate(to, subject, templateName string, data interface{}) (*Mail, error) {
	return NewMailBuilder().
		To(to).
		Subject(subject).
		Template(templateName, data).
		Build()
}

func (m *Mail) validate() error {
	if m.To == "" {
		return errors.New("recipient email cannot be empty")
	}
	if m.Subject == "" {
		return errors.New("subject cannot be empty")
	}
	if !m.IsHTML && m.Body == "" {
		return errors.New("plain text body cannot be empty")
	}
	if m.IsHTML && m.HTMLBody == "" {
		return errors.New("HTML body cannot be empty")
	}
	return nil
}

func (m *Mail) Send() error {
	config, err := NewSMTPConfig()
	if err != nil {
		return fmt.Errorf("SMTP configuration error: %w", err)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("SMTP configuration validation error: %w", err)
	}

	return m.SendWithConfig(config)
}

func (m *Mail) SendWithConfig(config *SMTPConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("SMTP configuration validation error: %w", err)
	}

	return m.sendWithConfig(config)
}

func (m *Mail) sendWithConfig(config *SMTPConfig) error {
	if err := m.validate(); err != nil {
		return fmt.Errorf("mail validation error: %w", err)
	}

	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	mail := gomail.NewMessage()
	mail.SetHeader("From", config.Username)
	mail.SetHeader("To", m.To)
	mail.SetHeader("Subject", m.Subject)

	if m.IsHTML {
		mail.SetBody("text/html", m.HTMLBody)
		if m.Body != "" {
			mail.AddAlternative("text/plain", m.Body)
		}
	} else {
		mail.SetBody("text/plain", m.Body)
	}

	if err := d.DialAndSend(mail); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
