package websocket

import (
	"log"
	"net/http"
	"strings"

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
	log.Println("Creating new WebSocket service")
	return &service{hub: hub}
}

func (s *service) HandleConnection(w http.ResponseWriter, r *http.Request, roomID string, userID string) {
	log.Printf("[WebSocket] HandleConnection - Room: %s, User: %s", roomID, userID)
	log.Printf("[WebSocket] Headers: %v", r.Header)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			allowed := origin == "http://localhost:5173" ||
				origin == "http://127.0.0.1:5173" ||
				origin == "http://localhost:8081" ||
				strings.HasPrefix(origin, "http://localhost:")
			if !allowed {
				log.Printf("[WebSocket] Origin not allowed: %s", origin)
			}
			return allowed
		},
		Subprotocols: []string{"Bearer"},
	}

	clientProtocols := websocket.Subprotocols(r)
	var selectedProtocol string

	for _, p := range clientProtocols {
		if p == "Bearer" || strings.HasPrefix(p, "Bearer ") {
			selectedProtocol = "Bearer"
			break
		}
	}

	conn, err := upgrader.Upgrade(w, r, http.Header{
		"Sec-WebSocket-Protocol": []string{selectedProtocol},
	})
	if err != nil {
		log.Printf("[WebSocket] Upgrade error: %v", err)
		return
	}

	log.Printf("[WebSocket] Connection upgraded successfully")

	client := NewClient(s.hub, conn, roomID, userID)
	log.Printf("[WebSocket] New client created: %+v", client)

	s.hub.register <- client
	log.Printf("[WebSocket] Client registered in hub")

	go client.writePump()
	go client.readPump()

	log.Printf("[WebSocket] Read/write pumps started")
}

func (s *service) CreateClient(conn *websocket.Conn, roomID, userID string) *Client {
	log.Printf("[WebSocket] Creating new client - Room: %s, User: %s", roomID, userID)
	return NewClient(s.hub, conn, roomID, userID)
}

func (s *service) Broadcast(roomID string, message Message) {
	log.Printf("Broadcasting message of type %s to room %s", message.Type, roomID)

	message.RoomID = roomID
	select {
	case s.hub.broadcast <- message:
		log.Printf("Message successfully queued for broadcast to room %s", roomID)
	default:
		log.Printf("Warning: Broadcast channel full for room %s, message dropped", roomID)
	}
}
