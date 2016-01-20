package web

import (
	"fmt"
	"galapagia/engine"
	"net/http"
)

func dataWebSocketHandler(gs *engine.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error during websockets upgrade:", err)
			return
		}

		for {
			var j map[string]interface{}
			err = conn.ReadJSON(&j)
			if err != nil {
				fmt.Println("Error reading JSON:", err)
				return
			}

			switch j["command"] {
			case "show current grid":
				go func() {
					j := gs.CurrentGrid()
					err = conn.WriteJSON(j)
					if err != nil {
						fmt.Println("Error writing JSON:", err)
						return
					}
				}()
			case "reset":
				go func() {
					gs.Reset()
					j := gs.CurrentGrid()
					err = conn.WriteJSON(j)
					if err != nil {
						fmt.Println("Error writing JSON:", err)
						return
					}
				}()
			default:
				fmt.Println("Unrecognized command:", j)
			}
		}
	}
}
