package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rahulbalajee/lenslocked/models"
)

// Decouple SessionService from controllers using interface
type SessionService interface {
	Create(userID int) (*models.Session, error)
	User(token string) (*models.User, error)
	Delete(token string) error
}

type Users struct {
	Templates struct {
		SignUp      Executer
		SignIn      Executer
		CurrentUser Executer
	}
	UserService    *models.UserService // tight coupling example (bad practice)
	SessionService SessionService      // decoupled with interface (best practice) Interface connection happens in line 46 in main.go
}

func (u Users) SignUp(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	u.Templates.SignUp.Execute(w, r, data)
}

func (u Users) ProcessSignUp(w http.ResponseWriter, r *http.Request) {
	// Get the email and password from the form
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Create a new user during signup
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Create a new session for the user
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: Show a warning message to the user
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	// Set the token returned from SessionService.Create in a Cookie in user's browser for authenticating future requests
	setCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	// Authenticate user using the email and check password with bcrypt
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	// TODO: Check for SQL ErrNoRows in case the user tries to login without signing up first
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Create a new session token for the user and set cookie
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		// Cookie doesn't exists, so redirect them to /signin page
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	// Pass the token to SessionService and get the *models.User
	user, err := u.SessionService.User(token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	u.Templates.CurrentUser.Execute(w, r, user)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	// Delete the session token from the DB and remove the cookie from user's browser = log them out
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	deleteCookie(w, CookieSession)

	http.Redirect(w, r, "/signin", http.StatusFound)
}
