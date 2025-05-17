package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan Message
	roomID string
	userID string
}

func NewClient(hub *Hub, conn *websocket.Conn, roomID, userID string) *Client {
	log.Printf("[WS Client] New client created - Room: %s, User: %s", roomID, userID)
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan Message, 256),
		roomID: roomID,
		userID: userID,
	}
}

func (c *Client) readPump() {
	log.Printf("[WS Client] Starting readPump for user %s in room %s", c.userID, c.roomID)
	defer func() {
		log.Printf("[WS Client] Closing readPump for user %s", c.userID)
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		log.Printf("[WS Client] Received pong from user %s", c.userID)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS Client] Unexpected close error: %v, User: %s", err, c.userID)
			} else {
				log.Printf("[WS Client] Connection closed for user %s", c.userID)
			}
			break
		}

		log.Printf("[WS Client] Raw message received from user %s: %s", c.userID, string(rawMsg))

		var msg Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Printf("[WS Client] Error unmarshaling message: %v, User: %s", err, c.userID)
			continue
		}

		msg.RoomID = c.roomID
		msg.UserID = c.userID
		log.Printf("[WS Client] Forwarding message of type %s from user %s", msg.Type, c.userID)
		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	log.Printf("[WS Client] Starting writePump for user %s", c.userID)
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Printf("[WS Client] Closing writePump for user %s", c.userID)
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Printf("[WS Client] Channel closed, sending close message to user %s", c.userID)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			msgBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("[WS Client] Error marshaling message: %v, User: %s", err, c.userID)
				continue
			}

			log.Printf("[WS Client] Sending message to user %s: %s", c.userID, string(msgBytes))
			if err := c.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				log.Printf("[WS Client] Error writing message: %v, User: %s", err, c.userID)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			log.Printf("[WS Client] Sending ping to user %s", c.userID)
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[WS Client] Error sending ping: %v, User: %s", err, c.userID)
				return
			}
		}
	}
}
