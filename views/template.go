package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Parse(filepath string) (Template, error) {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template %w", err)
	}
	return Template{htmltmpl: tmpl}, nil
}

type Template struct {
	htmltmpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmltmpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
