package websocket

import (
	"encoding/json"
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

func ValidateJWT(token string) (int, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/auth/validate", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("invalid token (status %d)", resp.StatusCode)
	}

	var response struct {
		UserID int `json:"userId"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response: %v", err)
	}

	return response.UserID, nil
}

func (s *service) Broadcast(roomID string, message Message) {
	message.RoomID = roomID
	s.hub.broadcast <- message
}
