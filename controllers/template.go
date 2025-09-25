package controllers

import "net/http"

type Executer interface {
	Execute(w http.ResponseWriter, r *http.Request, data any, errs ...error)
}
