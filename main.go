package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/rahulbalajee/lenslocked/controllers"
	"github.com/rahulbalajee/lenslocked/models"
	"github.com/rahulbalajee/lenslocked/templates"
	"github.com/rahulbalajee/lenslocked/views"
)

func main() {
	r := chi.NewRouter()

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

	cfg := models.DefaultPostgresConfig()
	// Init DB connection
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Dependency injection (passing in the PostgreSQL DB)
	userService := models.UserService{
		DB: db,
	}

	// Dependency injection (passing in the PostgreSQL DB)
	sessionService := models.SessionService{
		DB: db,
	}

	// Adapting REST and using it's own controllers for User related endpoints
	usersC := controllers.Users{
		UserService: &userService,
		// Interface connection for SessionService happens here
		SessionService: &sessionService,
	}
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
	r.Get("/signup", usersC.SignUp)
	r.Post("/signup", usersC.ProcessSignUp)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting server on :3000...")

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX" // TODO: Load this from an env var
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false), // TODO: Fix this before deploy
		csrf.TrustedOrigins([]string{"localhost:3000", "127.0.0.1:3000"}),
	)

	http.ListenAndServe(":3000", csrfMw(r))
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
