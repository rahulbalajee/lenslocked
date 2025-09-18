package main

import (
	"github.com/rahulbalajee/lenslocked/models"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "8dc083b1f8a082"
	password = "d9f5958e5b9ecc"
)

func main() {
	email := models.Email{
		From:      "test@lenslocked.com",
		To:        "jon@calhoun.io",
		Subject:   "This is a test email",
		Plaintext: "This is the body of the email",
		HTML:      `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`,
	}

	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	err := es.Send(email)
	if err != nil {
		panic(err)
	}

}
