package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan Message
	roomID string
	userID string
}

func NewClient(hub *Hub, conn *websocket.Conn, roomID, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan Message, 256),
		roomID: roomID,
		userID: userID,
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		msg.RoomID = c.roomID
		msg.UserID = c.userID
		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		msgBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("error marshaling message: %v", err)
			continue
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			log.Printf("error writing message: %v", err)
			break
		}
	}
}