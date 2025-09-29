package controllers

import (
	"net/http"
)

type Galleries struct {
	Template struct {
		New Executer
	}
	GalleryService GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Template.New.Execute(w, r, data)
}
