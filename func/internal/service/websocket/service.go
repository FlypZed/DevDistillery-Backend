package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Service interface {
	HandleConnection(w http.ResponseWriter, r *http.Request, roomID string, userID string)
	Broadcast(roomID string, message Message)
}

type service struct {
	hub *Hub
}

func NewWebSocketService(hub *Hub) Service {
	return &service{hub: hub}
}

func (s *service) HandleConnection(w http.ResponseWriter, r *http.Request, roomID string, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := NewClient(s.hub, conn, roomID, userID)
	s.hub.register <- client

	go client.writePump()
	go client.readPump()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:5173" || origin == "dev-distillery-produccuib.com" // Cambiar despues :D
	},
}

func ValidateJWT(token string) (string, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/auth/validate", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid token")
	}
	return "user-id-from-token", nil
}

func (s *service) Broadcast(roomID string, message Message) {
	message.RoomID = roomID
	s.hub.broadcast <- message
}
