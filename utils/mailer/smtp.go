package mailer

import (
	"errors"
	"senkou-catalyst-be/utils/config"
	"strconv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewSMTPConfig() (*SMTPConfig, error) {
	host := config.GetEnv("MAILER_HOST", "smtp.gmail.com")
	port := config.GetEnv("MAILER_PORT", "587")
	username := config.GetEnv("MAILER_USERNAME", "")
	password := config.GetEnv("MAILER_PASSWORD", "")

	if host == "" || port == "" || username == "" || password == "" {
		return nil, errors.New("SMTP configuration is incomplete")
	}

	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.New("invalid MAILER_PORT value")
	}

	return &SMTPConfig{
		Host:     host,
		Port:     parsedPort,
		Username: username,
		Password: password,
	}, nil
}

func (c *SMTPConfig) Validate() error {
	if c.Host == "" {
		return errors.New("SMTP host is required")
	}
	if c.Port <= 0 {
		return errors.New("SMTP port must be positive")
	}
	if c.Username == "" {
		return errors.New("SMTP username is required")
	}
	if c.Password == "" {
		return errors.New("SMTP password is required")
	}
	return nil
}
