package main

import (
	"http"
	"log"
	"websocket"
	"fmt"
)

var messageChan = make(chan []byte)
var subscriptionChan = make(chan subscription)

type subscription struct {
	conn      *websocket.Conn
	subscribe bool
}

func main() {
	go hub()

	http.HandleFunc("/ws", webSocketProtocolSwitch)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func hub() {
	conns := make(map[*websocket.Conn]int)
	for {
		select {
			case subscription := <-subscriptionChan:
				fmt.Printf("subscription: %v", subscription)
				conns[subscription.conn] = 0, subscription.subscribe
			case message := <-messageChan:
				fmt.Printf("message: %v", message)
				for conn, _ := range conns {
					if _, err := conn.Write(message); err != nil {
						conn.Close()
					}
				}
		}
	}
}

func webSocketProtocolSwitch(c http.ResponseWriter, req *http.Request) {
	// Handle old and new versions of protocol.
	if _, found := req.Header["Sec-Websocket-Key1"]; found {
		websocket.Handler(clientHandler).ServeHTTP(c, req)
	} else {
		websocket.Draft75Handler(clientHandler).ServeHTTP(c, req)
	}
}

func clientHandler(ws *websocket.Conn) {
	defer func() {
		subscriptionChan <- subscription{ws, false}
		ws.Close()
	}()

	subscriptionChan <- subscription{ws, true}

	buf := make([]byte, 256)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			break
		}
		messageChan <- buf[0:n]
	}
}
