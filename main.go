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

	tmpl, err := views.Parse("templates/home.gohtml")
	if err != nil {
		panic(err)
	}
	r.Get("/", controllers.StaticHandler(tmpl))

	tmpl, err = views.Parse("templates/contact.gohtml")
	if err != nil {
		panic(err)
	}
	r.Get("/contact", controllers.StaticHandler(tmpl))

	tmpl, err = views.Parse("templates/faq.gohtml")
	if err != nil {
		panic(err)
	}
	r.Get("/faq", controllers.StaticHandler(tmpl))

	//r.Use(middleware.Logger)
	//r.Get("/page/{page-id}", newPageHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting server on :3000...")
	http.ListenAndServe(":3000", r)
}
