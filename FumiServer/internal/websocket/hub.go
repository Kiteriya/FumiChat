package websocket

import (
	"FumiServer/internal/models"
	"encoding/json"
	"log"
)

type Hub struct {
	route      chan models.Message
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		route:      make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func Run(h *Hub) {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
			log.Println("User connected")

		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
				log.Println("User disconnected")
			}

		case message := <-h.route:
			sendBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Cant marshal json: %v", err)
				continue
			}

			targetClient, ok := h.clients[message.TargetID]
			if ok != false {
				log.Printf("Failed to find client: %v", ok)
				continue
			}

			targetClient.send <- sendBytes
		}

	}
}
