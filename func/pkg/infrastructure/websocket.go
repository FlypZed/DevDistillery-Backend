package infrastructure

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

type BoardMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	for {
		var msg BoardMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Error reading message:", err)
			break
		}

		handleMessage(msg, conn)
	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
}

func handleMessage(msg BoardMessage, conn *websocket.Conn) {
	switch msg.Type {
	case "draw":
		broadcast(msg)
	case "shape":
		broadcast(msg)
	default:
		log.Println("Unknown message type:", msg.Type)
	}
}

func broadcast(msg BoardMessage) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Println("Error broadcasting message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
