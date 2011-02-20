package main

import (
	"http"
	"log"
	"websocket"
	"time"
	"math"
	"strconv"
)

func main() {
	log.Println("Starting Server...")

	http.Handle("/ws", websocket.Handler(handler));

	err := http.ListenAndServe(":8080", nil);

	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}

func handler(ws *websocket.Conn) {
	defer func() {
		log.Printf("Closing websocket: %v\n", ws)
        	ws.Close()
	}()

	x := 0.
	for {
		if x >= 2*math.Pi {
			x = 0
		} else {
			x += 0.05
		}
		
		time.Sleep(500*1000*1000) // sleep for 500ms (Sleep takes nanoseconds)

		msg := strconv.Ftoa64(math.Sin(x), 'g', 10)
		log.Printf("%v sending: %v\n", ws, msg)
		ws.Write( []byte(msg) )
	}
}

