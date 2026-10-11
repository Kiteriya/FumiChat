package main

import (
	auth "FumiServer/internal/autorization"
	"FumiServer/internal/database"
	"FumiServer/internal/models"
	ws "FumiServer/internal/websocket"
	"encoding/json"

	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}
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
		token, err := auth.LoginHandler(database.DB, w, r)
		if err != nil {
			return
		}

		response := models.AuthResponse{
			JwtTOKEN: token,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		respError := json.NewEncoder(w).Encode(response)
		if respError != nil {
			log.Printf("Ошибка при отправке JSON: %v", err)
		}
		return
	})

	log.Println("Server is running")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
