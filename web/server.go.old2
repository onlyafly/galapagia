package web

import (
	"fmt"
	"html/template"
	"net/http"

	"galapagia/Godeps/_workspace/src/github.com/gorilla/mux"
	"galapagia/Godeps/_workspace/src/github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func print_binary(s []byte) {
	fmt.Printf("Received b:")
	for n := 0; n < len(s); n++ {
		fmt.Printf("%d,", s[n])
	}
	fmt.Printf("\n")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		print_binary(p)

		err = conn.WriteMessage(messageType, p)
		if err != nil {
			return
		}
	}
}

func handleDataWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		fmt.Println("Waiting to read data over websocket")

		/*
			_, p, _ := conn.ReadMessage()
			fmt.Println("Read from ReadMessage", p)*/

		var j map[string]interface{}
		err = conn.ReadJSON(&j)
		if err != nil {
			fmt.Println("Error reading JSON:", err)
			return
		}
		fmt.Println("Read over the websocket:", j)

		c := CounterJson{
			Counter: 64,
		}
		err = conn.WriteJSON(c)
		if err != nil {
			fmt.Println("Error writing JSON:", err)
			return
		}
	}
}

func ServeSite() {
	router := mux.NewRouter().StrictSlash(true)
	templates = template.Must(template.New("").ParseGlob("web/templates/*")) // compile all templates and cache them

	router.Methods("GET").Path("/echo").HandlerFunc(echoHandler)
	router.Methods("GET").Path("/").HandlerFunc(homeHandler)
	router.Methods("GET").Path("/api/dataWebSocket").HandlerFunc(handleDataWebSocket)
	router.Methods("GET").Path("/api/tick").HandlerFunc(handleGetTick)

	// Serve all static files
	// ORDERING: Must be after all other routes
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web/public/"))))

	fmt.Println("Serving galapagia on :8080...")
	http.ListenAndServe(":8080", router)
}
