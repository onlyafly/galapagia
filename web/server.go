package web

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	templates *template.Template
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	c := map[string]string{
		"PageTitle": "galapagia",
	}

	respondWithTemplate(templates, "index", c, w, r)
}

func ServeSite() {
	router := mux.NewRouter().StrictSlash(true)
	templates = template.Must(template.New("").ParseGlob("web/templates/*")) // compile all templates and cache them

	router.Methods("GET").Path("/").HandlerFunc(homeHandler)

	// Serve all static files
	// ORDERING: Must be after all other routes
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web/public/"))))

	http.ListenAndServe(":8080", router)
}
