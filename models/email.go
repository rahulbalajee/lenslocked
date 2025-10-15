package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.bc-noc-dev.com"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type EmailService struct {
	//  We can also add a DefaultSender field that can be set if needed, otherwise we will use a constant defined in our code
	DefaultSender string

	// unexported field because the caller doesn't need to know about our implementation
	dialer *mail.Dialer
}

// Factory to construct EmailService and the unexported field
func NewEmailService(config SMTPConfig) *EmailService {
	return &EmailService{
		dialer: mail.NewDialer(
			config.Host,
			config.Port,
			config.Username,
			config.Password),
	}
}

func (es *EmailService) Send(email Email) error {
	msg := mail.NewMessage()

	// setFrom() to a default value in case it's not set in Email
	es.setFrom(msg, email)

	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)

	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}

	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func (es *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		Subject:   "Reset your password",
		To:        to,
		Plaintext: "To reset your password, please visit the following link: " + resetURL,
		HTML:      `<p>To reset your password, please visit the following link: <a href="` + resetURL + `">` + resetURL + `</a></p>`,
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}

	return nil
}

func (es *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string

	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}

	msg.SetHeader("From", from)
}
