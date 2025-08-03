package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email us at <a href=\"mailto:rahrkb4@gmail.com\">rahrkb4@gmail.com</a>.")
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

type Router struct{}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		http.NotFound(w, r)
		//http.Error(w, fmt.Sprint(errors.New("no page found")), http.StatusNotFound)
	}
}

func main() {
	var router *Router
	fmt.Println("Starting server on :3000...")
	http.ListenAndServe(":3000", router)
}
