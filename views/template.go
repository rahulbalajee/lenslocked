package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
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
	tmpl := template.New(patterns[0])

	tmpl = tmpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return `<!-- TODO: Implement the CSRFField -->`
			},
		},
	)

	tmpl, err := tmpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template from FS %w", err)
	}

	return Template{htmltmpl: tmpl}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	tmpl, err := t.htmltmpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}
	tmpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
