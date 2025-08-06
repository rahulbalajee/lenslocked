package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("executing template %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/home.gohtml")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/contact.gohtml")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w,
		`
	<h1>FAQ Page</h1>
	<ul>
		<li><b>Is there a free version?</b>Yeah G!</li>
		<li>What are your support hours? 9 to 5 G!</li>
	</ul>
	`)
}

func newPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	page := chi.URLParamFromCtx(ctx, "id")

	fmt.Fprintf(w, "The page ID is %s", page)
}

func main() {
	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.Route("/galleries", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/{id}", newPageHandler)
	})
	//r.Get("/page/{page-id}", newPageHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting server on :3000...")
	http.ListenAndServe(":3000", r)
}
