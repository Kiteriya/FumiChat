package main

import (
	auth "FumiServer/internal/autorization"
	"FumiServer/internal/database"
	ws "FumiServer/internal/websocket"

	"log"
	"net/http"
)

func main() {
	database.InitDatabase()

	hub := ws.NewHub()
	go ws.Run(hub)

	http.HandleFunc("/FumiServer/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ConnectionHandeler(hub, w, r)
	})

	http.HandleFunc("/FumiServer/register", func(w http.ResponseWriter, r *http.Request) {
		auth.RegisterHandler(database.DB, w, r)
	})

	http.HandleFunc("/FumiServer/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginHandler(w, r)
	})

	log.Println("Server is running")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
