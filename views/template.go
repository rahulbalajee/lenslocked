package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/csrf"
	"github.com/rahulbalajee/lenslocked/context/context"
	"github.com/rahulbalajee/lenslocked/models"
)

type public interface {
	Public() string
}

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
	// We need this to be able to add custom function to our templates
	tmpl := template.New(path.Base(patterns[0]))

	// We need to add the csrfField function before parsing the template or will get an error
	// No access to request r here so just stubbing
	// Template creation: You register stub functions just so the template parser doesn't crash when it sees {{csrfField}} in your .gohtml files when parsing
	tmpl = tmpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("currentUser not implemented")
			},
			"errors": func() []string {
				return nil
			},
		},
	)

	// After adding custom function using FuncMap, parse the templates
	tmpl, err := tmpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template from FS %w", err)
	}

	return Template{htmltmpl: tmpl}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any, errs ...error) {
	// Templates can be used concurrently by multiple goroutines (multiple users hitting your site at once)
	// If you modified the original template directly, you'd have race conditions where one user's request data might leak into another user's response
	// Cloning gives each request its own isolated copy with the correct request-specific data
	tmpl, err := t.htmltmpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}

	errMsgs := errMessages(errs...)

	tmpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}

func errMessages(errs ...error) []string {
	var msgs []string

	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}

	return msgs
}
