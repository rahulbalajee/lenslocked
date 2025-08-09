package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rahulbalajee/lenslocked/controllers"
	"github.com/rahulbalajee/lenslocked/views"
)

func main() {
	r := chi.NewRouter()

	tmpl := views.Must(views.Parse("templates/home.gohtml"))
	r.Get("/", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.Parse("templates/contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.Parse("templates/faq.gohtml"))
	r.Get("/faq", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.Parse("templates/newpage.gohtml"))
	r.Get("/newpage", controllers.StaticHandler(tmpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting server on :3000...")
	http.ListenAndServe(":3000", r)
}
