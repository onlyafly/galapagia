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

type DataWebSocket struct {
	Handler    http.HandlerFunc
	FromClient chan map[string]interface{}
	ToClient   chan interface{}
}

func NewDataWebSocket() *DataWebSocket {
	from := make(chan map[string]interface{})
	to := make(chan interface{})

	h := func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		go func() {
			for {
				j := <-to
				err = conn.WriteJSON(j)
				if err != nil {
					fmt.Println("Error writing JSON:", err)
					return
				}
			}
		}()

		for {
			var j map[string]interface{}
			err = conn.ReadJSON(&j)
			if err != nil {
				fmt.Println("Error reading JSON:", err)
				return
			} else {
				from <- j
			}
		}
	}

	return &DataWebSocket{
		FromClient: from,
		ToClient:   to,
		Handler:    h,
	}
}

func ServeSite() {
	router := mux.NewRouter().StrictSlash(true)
	templates = template.Must(template.New("").ParseGlob("web/templates/*")) // compile all templates and cache them

	router.Methods("GET").Path("/echo").HandlerFunc(echoHandler)
	router.Methods("GET").Path("/").HandlerFunc(homeHandler)

	dataWebSocket := NewDataWebSocket()

	router.Methods("GET").Path("/api/dataWebSocket").HandlerFunc(dataWebSocket.Handler)
	go func() {
		for {
			fmt.Println("Waiting for a read from client:")
			m := <-dataWebSocket.FromClient
			fmt.Println("Read over the websocket:", m)
		}
	}()
	go func() {
		c := CounterJson{
			Counter: 64,
		}
		fmt.Println("Sending data to client:", c)
		dataWebSocket.ToClient <- c
		fmt.Println("Sent data to client:", c)
	}()

	router.Methods("GET").Path("/api/tick").HandlerFunc(handleGetTick)

	// Serve all static files
	// ORDERING: Must be after all other routes
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web/public/"))))

	fmt.Println("Serving galapagia on :8080...")
	http.ListenAndServe(":8080", router)
}
