package main

import (
	"fmt"
	"os"

	"github.com/go-mail/mail/v2"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "8dc083b1f8a082"
	password = "d9f5958e5b9ecc"
)

func main() {
	from := "test@lenslocked.com"
	to := "jon@calhoun.io"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)
	msg.WriteTo(os.Stdout)

	dialer := mail.NewDialer(host, port, username, password)
	err := dialer.DialAndSend(msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Message sent")

}

/*

MIME-Version: 1.0
Date: Wed, 17 Sep 2025 23:20:13 +0530
From: test@lenslocked.com
Subject: This is a test email
To: jon@calhoun.io
Content-Type: multipart/alternative;
 boundary=fb0e850fcfd27eb855b667d000746610d45902bf460631ba5614845b97f9

--fb0e850fcfd27eb855b667d000746610d45902bf460631ba5614845b97f9
Content-Transfer-Encoding: quoted-printable
Content-Type: text/plain; charset=UTF-8

This is the body of the email
--fb0e850fcfd27eb855b667d000746610d45902bf460631ba5614845b97f9
Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=UTF-8

<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>
--fb0e850fcfd27eb855b667d000746610d45902bf460631ba5614845b97f9--

*/
