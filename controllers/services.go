package controllers

import (
	"io"

	"github.com/rahulbalajee/lenslocked/models"
)

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

type GalleryService interface {
	Create(title string, userID int) (*models.Gallery, error)
	ByID(id int) (*models.Gallery, error)
	ByUserID(userID int) ([]models.Gallery, error)
	Update(gallery *models.Gallery) error
	Delete(id int) error
	Images(galleryID int) ([]models.Image, error)
	Image(galleryId int, filename string) (models.Image, error)
	DeleteImage(galleryID int, filename string) error
	CreateImage(galleryID int, filename string, contents io.ReadSeeker) error
}
