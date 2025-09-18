package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rahulbalajee/lenslocked/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	host := os.Getenv("SMTP_HOST")

	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	err = es.ForgotPassword("jon@calhoun.io", "http://sampleurl")
	if err != nil {
		panic(err)
	}

}
