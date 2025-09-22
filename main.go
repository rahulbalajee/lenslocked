package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/rahulbalajee/lenslocked/controllers"
	"github.com/rahulbalajee/lenslocked/migrations"
	"github.com/rahulbalajee/lenslocked/models"
	"github.com/rahulbalajee/lenslocked/templates"
	"github.com/rahulbalajee/lenslocked/views"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	// TODO: read Postgres config from env
	cfg.PSQL = models.DefaultPostgresConfig()

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure, err = strconv.ParseBool(os.Getenv("CSRF_SECURE"))
	if err != nil {
		return cfg, err
	}

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Initiate DB connection and close it later when function exits
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Run DB migrations automatically at startup
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Dependency injection (passing in the PostgreSQL DB)
	//userService for creating and managing users
	userService := &models.UserService{
		DB: db,
	}

	// Dependency injection (passing in the PostgreSQL DB)
	//sessionService for creating and managing session
	sessionService := &models.SessionService{
		DB: db,
	}

	passwordResetService := &models.PasswordResetService{
		DB: db,
	}

	emailService := models.NewEmailService(cfg.SMTP)

	// Setup our middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.TrustedOrigins([]string{"localhost:3000", "127.0.0.1:3000"}),
	)

	// Adapting REST and using it's own controllers for User related endpoints plumbing UserService and SessionService
	usersC := controllers.Users{
		UserService: userService,
		// Interface connection for SessionService happens here (plumbing done)
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}

	// Plumbing work to make sure the Templates in users controller are populated with right values before routing happens
	usersC.Templates.SignUp = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml",
		"tailwind.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml",
		"tailwind.gohtml",
	))
	usersC.Templates.CurrentUser = views.Must(views.ParseFS(
		templates.FS,
		"current-user.gohtml",
		"tailwind.gohtml",
	))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-pw.gohtml",
		"tailwind.gohtml",
	))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS,
		"check-your-email.gohtml",
		"tailwind.gohtml",
	))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS,
		"reset-password.gohtml",
		"tailwind.gohtml",
	))

	// Setup our router
	r := chi.NewRouter()

	// Apply middlewares
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	tmpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	// Static handler expects controllers.Executor type, but we can pass in Views.Template
	// because both types implement the exact same method Execute
	// Execute(w http.ResponseWriter, r *http.Request, data any)
	// This is where the interface connection happens
	// With the Executor interface we're decoupling controllers package from views package
	r.Get("/", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tmpl))

	// Routing happens here
	r.Get("/signup", usersC.SignUp)
	r.Post("/signup", usersC.ProcessSignUp)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	//r.Post("/signout", usersC.ProcessSignOut)
	r.With(umw.RequireUser).Post("/signout", usersC.ProcessSignOut)

	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	//r.Get("/users/me", usersC.CurrentUser)
	// Create a subrouter for "/users/me" with RequireUser middleware
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.With(umw.RequireUser).Post("/update-email", usersC.UpdateEmail)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Printf("Starting server on %s...", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
