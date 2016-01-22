package web

import (
	"fmt"
	"galapagia/engine"
	"net/http"
	"time"

	"galapagia/Godeps/_workspace/src/github.com/gorilla/websocket"
)

func generateStateAccessChannels() (commandChan chan<- string, dataChan <-chan [][]int) {
	commands := make(chan string)
	data := make(chan [][]int)
	s := engine.NewState(100, 100)

	go func() {
		for {
			m := <-commands
			switch m {
			case "log":
				s.LogBugs()
			case "tick":
				s.Tick()
			case "reset":
				s.Reset(100)
			case "grid":
				data <- s.CurrentCellGrid()
			default:
				fmt.Println("Unrecognized state access command", m)
			}
		}
	}()

	return commands, data
}

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

func dataWebSocketHandler() http.HandlerFunc {
	commandChan, dataChan := generateStateAccessChannels()

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
					commandChan <- "grid"
					writeChan <- <-dataChan
				}()
			case "reset":
				go func() {
					commandChan <- "reset"
					commandChan <- "grid"
					writeChan <- <-dataChan
				}()
			case "tick 30 times":
				go func() {
					for i := 0; i < 3000; i++ {
						commandChan <- "tick"
						commandChan <- "grid"
						writeChan <- <-dataChan
						time.Sleep(10 * time.Millisecond)
					}
				}()
			case "log data":
				go func() {
					commandChan <- "log"
				}()
			default:
				fmt.Println("Unrecognized command:", j)
			}
		}
	}
}
