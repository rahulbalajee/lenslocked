package models

import "github.com/go-mail/mail/v2"

const (
	DefaultSender = "support@lenslocked.com"
)

type EmailService struct {
	DefaultSender string

	// unexported field because the caller doesn't need to know about our implementation
	dialer *mail.Dialer
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Factory to construct EmailService and the unexported field
func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(config.Host,
			config.Port,
			config.Username,
			config.Password),
	}

	return &es
}
