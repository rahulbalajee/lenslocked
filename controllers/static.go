package controllers

import (
	"html/template"
	"net/http"
)

// Closure that captures the template and returns http.HandlerFunc required by Chi router
func StaticHandler(tmpl Executer) http.HandlerFunc {
	// This anonymous function is the closure
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, r, nil) // 'tmpl' is captured from outer scope
	}
}

func FAQ(tmpl Executer) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! We offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "What are your support hours?",
			Answer:   "We have support staff answering emails 24/7, though response times may be a bit slower on weekends.",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>`,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, r, questions)
	}
}
