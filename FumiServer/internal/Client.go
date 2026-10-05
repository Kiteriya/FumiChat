package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Message struct {
	messageType string `json:"message_type"`
	senderId    string `json:"sender_id"`
	targetId    string `json:"target_id"`
	message     string `json:"payload"`
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}
		msg.senderId = c.Id
		c.hub.route <- msg
	}
}

func (c *Client) WritePump() {

}
