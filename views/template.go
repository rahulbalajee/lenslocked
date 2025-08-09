package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmltmpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tmpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template from FS %w", err)
	}
	return Template{htmltmpl: tmpl}, nil
}

/*
func Parse(filepath string) (Template, error) {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template %w", err)
	}
	return Template{htmltmpl: tmpl}, nil
}
*/

func (t Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmltmpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
