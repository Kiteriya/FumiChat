package main

import (
	websocket "FumiServer/internal"
	"log"
	"net/http"
)

func main() {
	hub := websocket.NewHub()
	go websocket.Run(hub)

	http.HandleFunc("/FumiServer", func(w http.ResponseWriter, r *http.Request) {
		websocket.ConnectionHandeler(hub, w, r)
	})

	log.Println("Server is running")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
