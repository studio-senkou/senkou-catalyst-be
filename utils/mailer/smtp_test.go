package mailer

import (
	"os"
	"testing"
)

func TestNewSMTPConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"MAILER_HOST":     "smtp.example.com",
				"MAILER_PORT":     "587",
				"MAILER_USERNAME": "user@example.com",
				"MAILER_PASSWORD": "password123",
			},
			expectError: false,
		},
		{
			name: "missing host",
			envVars: map[string]string{
				"MAILER_HOST":     "",
				"MAILER_PORT":     "587",
				"MAILER_USERNAME": "user@example.com",
				"MAILER_PASSWORD": "password123",
			},
			expectError: true,
			errorMsg:    "SMTP configuration is incomplete",
		},
		{
			name: "missing port",
			envVars: map[string]string{
				"MAILER_HOST":     "smtp.example.com",
				"MAILER_PORT":     "",
				"MAILER_USERNAME": "user@example.com",
				"MAILER_PASSWORD": "password123",
			},
			expectError: true,
			errorMsg:    "SMTP configuration is incomplete",
		},
		{
			name: "missing username",
			envVars: map[string]string{
				"MAILER_HOST":     "smtp.example.com",
				"MAILER_PORT":     "587",
				"MAILER_USERNAME": "",
				"MAILER_PASSWORD": "password123",
			},
			expectError: true,
			errorMsg:    "SMTP configuration is incomplete",
		},
		{
			name: "missing password",
			envVars: map[string]string{
				"MAILER_HOST":     "smtp.example.com",
				"MAILER_PORT":     "587",
				"MAILER_USERNAME": "user@example.com",
				"MAILER_PASSWORD": "",
			},
			expectError: true,
			errorMsg:    "SMTP configuration is incomplete",
		},
		{
			name: "invalid port format",
			envVars: map[string]string{
				"MAILER_HOST":     "smtp.example.com",
				"MAILER_PORT":     "invalid",
				"MAILER_USERNAME": "user@example.com",
				"MAILER_PASSWORD": "password123",
			},
			expectError: true,
			errorMsg:    "invalid MAILER_PORT value",
		},
		{
			name: "default values with valid override",
			envVars: map[string]string{
				"MAILER_USERNAME": "user@example.com",
				"MAILER_PASSWORD": "password123",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("MAILER_HOST")
			os.Unsetenv("MAILER_PORT")
			os.Unsetenv("MAILER_USERNAME")
			os.Unsetenv("MAILER_PASSWORD")

			for key, value := range tt.envVars {
				if value != "" {
					os.Setenv(key, value)
				}
			}

			config, err := NewSMTPConfig()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if config == nil {
					t.Errorf("expected config but got nil")
				}
			}
		})
	}

	os.Unsetenv("MAILER_HOST")
	os.Unsetenv("MAILER_PORT")
	os.Unsetenv("MAILER_USERNAME")
	os.Unsetenv("MAILER_PASSWORD")
}

func TestSMTPConfig_Validate(t *testing.T) {
	tests := []struct {
		name     string
		config   *SMTPConfig
		errorMsg string
	}{
		{
			name: "valid config",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "user@example.com",
				Password: "password123",
			},
			errorMsg: "",
		},
		{
			name: "empty host",
			config: &SMTPConfig{
				Host:     "",
				Port:     587,
				Username: "user@example.com",
				Password: "password123",
			},
			errorMsg: "SMTP host is required",
		},
		{
			name: "zero port",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Port:     0,
				Username: "user@example.com",
				Password: "password123",
			},
			errorMsg: "SMTP port must be positive",
		},
		{
			name: "negative port",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Port:     -1,
				Username: "user@example.com",
				Password: "password123",
			},
			errorMsg: "SMTP port must be positive",
		},
		{
			name: "empty username",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "",
				Password: "password123",
			},
			errorMsg: "SMTP username is required",
		},
		{
			name: "empty password",
			config: &SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "user@example.com",
				Password: "",
			},
			errorMsg: "SMTP password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.errorMsg == "" {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error '%s' but got none", tt.errorMsg)
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			}
		})
	}
}
