package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	templates *template.Template
)

func RespondWithJson(o interface{}, w http.ResponseWriter) error {
	js, err := json.Marshal(o)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func RespondWithTemplate(
	templateName string, templateContext interface{}, w http.ResponseWriter,
	r *http.Request) error {

	// you access the cached templates with the defined name, not the filename
	return templates.ExecuteTemplate(w, templateName, templateContext)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	c := map[string]string{
		"PageTitle": "galapagia",
	}

	RespondWithTemplate("index", c, w, r)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	templates = template.Must(template.New("").ParseGlob("web/templates/*")) // compile all templates and cache them

	router.Methods("GET").Path("/").HandlerFunc(homeHandler)

	// Serve all static files
	// ORDERING: Must be after all other routes
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web/public/"))))

	http.ListenAndServe(":8080", router)
}
