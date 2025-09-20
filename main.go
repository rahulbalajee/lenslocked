package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/rahulbalajee/lenslocked/controllers"
	"github.com/rahulbalajee/lenslocked/migrations"
	"github.com/rahulbalajee/lenslocked/models"
	"github.com/rahulbalajee/lenslocked/templates"
	"github.com/rahulbalajee/lenslocked/views"
)

func main() {
	// Load default config TODO: Fix in production before deploy
	cfg := models.DefaultPostgresConfig()

	// Init DB connection
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Run DB migrations
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Dependency injection (passing in the PostgreSQL DB)
	//userService for creating and managing users
	userService := models.UserService{
		DB: db,
	}

	// Dependency injection (passing in the PostgreSQL DB)
	//sessionService for creating and managing session
	sessionService := models.SessionService{
		DB: db,
	}

	// Setup our middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX" // TODO: Load this from an env var before production deploy
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false), // TODO: Fix this before deploy
		csrf.TrustedOrigins([]string{"localhost:3000", "127.0.0.1:3000"}),
	)

	// Adapting REST and using it's own controllers for User related endpoints plumbing UserService and SessionService
	usersC := controllers.Users{
		UserService: &userService,
		// Interface connection for SessionService happens here (plumbing done)
		SessionService: &sessionService,
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

	//r.Get("/users/me", usersC.CurrentUser)
	// Create a subrouter for "/users/me" with RequireUser middleware
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Starting server on :3000...")
	http.ListenAndServe(":3000", r)
}

/*
func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request came in from %s to %s\n", r.RemoteAddr, r.RequestURI)
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("time taken to serve request %v\n", time.Since(start))
	})
}
*/
