package web

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func respondWithJson(o interface{}, w http.ResponseWriter) error {
	js, err := json.Marshal(o)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func respondWithTemplate(
	templates *template.Template, templateName string, templateContext interface{}, w http.ResponseWriter,
	r *http.Request) error {

	// you access the cached templates with the defined name, not the filename
	return templates.ExecuteTemplate(w, templateName, templateContext)
}
