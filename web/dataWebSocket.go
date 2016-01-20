package web

import (
	"fmt"
	"net/http"
)

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
