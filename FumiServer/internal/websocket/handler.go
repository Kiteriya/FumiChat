package websocket

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectionHandler(h *Hub, w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer")

	secretKey := os.Getenv("JWT_SECRET")
	ok, userId := validateToken(tokenString, []byte(secretKey))
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{
		id:   userId,
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256)}

	client.hub.register <- client

	go client.ReadPump()
	go client.WritePump()
}

func validateToken(tokenString string, secretKey []byte) (bool, string) {

	log.Printf("Try parse token")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil || token == nil || !token.Valid {
		log.Printf("Token validation failed: %v", err)
		return false, ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Claims is not ok")
		return false, ""
	}

	userId, ok := claims["jti"].(string)
	if !ok {
		log.Printf("userId is not ok: %v", userId)
		return false, ""
	}

	log.Printf("Return err")
	return err == nil && token.Valid, userId
}
