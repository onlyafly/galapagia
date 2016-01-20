package web

import (
	"fmt"
	"galapagia/engine"
	"net/http"
	"time"

	"galapagia/Godeps/_workspace/src/github.com/gorilla/websocket"
)

/*
Connections support one concurrent reader and one concurrent writer.
Applications are responsible for ensuring that no more than one goroutine calls
the write methods (NextWriter, SetWriteDeadline, WriteMessage, WriteJSON)
concurrently and that no more than one goroutine calls the read methods
(NextReader, SetReadDeadline, ReadMessage, ReadJSON, SetPongHandler,
SetPingHandler) concurrently.
 -- http://www.gorillatoolkit.org/pkg/websocket
*/
func generateConnWriter(conn *websocket.Conn) chan<- interface{} {
	c := make(chan interface{})
	go func() {
		for {
			j := <-c
			err := conn.WriteJSON(j)
			if err != nil {
				fmt.Println("Error writing JSON:", err)
				return
			}
		}
	}()
	return c
}

func dataWebSocketHandler(gs *engine.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error during websockets upgrade:", err)
			return
		}

		writeChan := generateConnWriter(conn)

		for {
			var j map[string]interface{}
			err = conn.ReadJSON(&j)
			if err != nil {
				fmt.Println("Error reading JSON:", err)
				return
			}

			fmt.Println("Executing command:", j["command"])

			switch j["command"] {
			case "show current grid":
				go func() {
					writeChan <- gs.CurrentCellGrid()
				}()
			case "reset":
				go func() {
					gs.Reset()
					writeChan <- gs.CurrentCellGrid()
				}()
			case "tick 30 times":
				go func() {
					for i := 0; i < 30; i++ {
						gs.Tick()
						writeChan <- gs.CurrentCellGrid()
						time.Sleep(500 * time.Millisecond)
					}
				}()
			default:
				fmt.Println("Unrecognized command:", j)
			}
		}
	}
}
