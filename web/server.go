package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	templates *template.Template
	counter   int
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	c := map[string]string{
		"PageTitle": "alapagia",
	}

	respondWithTemplate(templates, "index", c, w, r)
}

type CounterJson struct {
	Counter int `json:"counter"`
}

func handleGetTick(w http.ResponseWriter, r *http.Request) {
	counter++

	c := CounterJson{
		Counter: counter,
	}

	respondWithJson(c, w)
}

func ServeSite() {
	router := mux.NewRouter().StrictSlash(true)
	templates = template.Must(template.New("").ParseGlob("web/templates/*")) // compile all templates and cache them

	router.Methods("GET").Path("/").HandlerFunc(homeHandler)
	router.Methods("GET").Path("/api/tick").HandlerFunc(handleGetTick)

	// Serve all static files
	// ORDERING: Must be after all other routes
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web/public/"))))

	fmt.Println("Serving galapagia on :8080...")
	http.ListenAndServe(":8080", router)
}
