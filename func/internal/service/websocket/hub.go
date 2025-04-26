package websocket

import (
	"log"
	"sync"
)

type Hub struct {
	rooms      map[string]*Room
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			room, exists := h.rooms[client.roomID]
			if !exists {
				room = NewRoom(client.roomID)
				h.rooms[client.roomID] = room
			}
			room.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered to room %s", client.roomID)

		case client := <-h.unregister:
			h.mu.Lock()
			if room, exists := h.rooms[client.roomID]; exists {
				if _, ok := room.clients[client]; ok {
					delete(room.clients, client)
					close(client.send)
					if len(room.clients) == 0 {
						delete(h.rooms, client.roomID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client unregistered from room %s", client.roomID)

		case message := <-h.broadcast:
			h.mu.RLock()
			if room, exists := h.rooms[message.RoomID]; exists {
				for client := range room.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(room.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}
