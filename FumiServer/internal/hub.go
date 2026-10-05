package websocket

import (
	"log"
)

type Hub struct {
	route      chan Message
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		route:      make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func Run(h *Hub) {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Id] = client
			log.Println("User connected")

		case client := <-h.unregister:
			if _, ok := h.clients[client.Id]; ok {
				delete(h.clients, client.Id)
				close(client.send)
				log.Println("User disconnected")
			}
		}

	}
}
