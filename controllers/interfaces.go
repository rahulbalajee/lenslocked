package controllers

import "github.com/rahulbalajee/lenslocked/models"

// Decouple SessionService from controllers using interface
type SessionService interface {
	Create(userID int) (*models.Session, error) // Need to decouple this completely by defining our own types
	User(token string) (*models.User, error)
	Delete(token string) error
}

// Decouple PasswordResetService from controllers using interfaces
type PasswordResetService interface {
	Create(email string) (*models.PasswordReset, error)
	Consume(token string) (*models.User, error)
}

type EmailService interface {
	ForgotPassword(to string, resetURL string) error
	Send(email models.Email) error
}
