package controllers

import "net/http"

type Executer interface {
	Execute(w http.ResponseWriter, data any)
}
