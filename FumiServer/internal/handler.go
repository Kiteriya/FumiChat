package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectionHandeler(h *Hub, w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	userId := "" //must be upgrade later

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{
		Id:   userId,
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256)}

	client.hub.register <- client

	go client.ReadPump()
	go client.WritePump()
}
