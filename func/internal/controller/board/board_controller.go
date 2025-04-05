package board

import (
	"context"
	"net/http"
	"sync"

	"func/internal/domain"
	service "func/internal/service/board"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type BoardMessage struct {
	Type  string         `json:"type"`
	Nodes []domain.Node  `json:"nodes,omitempty"`
	Edges []domain.Edge  `json:"edges,omitempty"`
	Data  map[string]any `json:"data,omitempty"`
}

type BoardController struct {
	service service.Service
	hubs    map[string]*BoardHub
	mutex   sync.Mutex
}

func NewBoardController(service service.Service) *BoardController {
	return &BoardController{
		service: service,
		hubs:    make(map[string]*BoardHub),
	}
}

type Client struct {
	hub  *BoardHub
	conn *websocket.Conn
	send chan BoardMessage
}

type BoardHub struct {
	boardID   string
	service   service.Service
	clients   map[*Client]bool
	register  chan *Client
	broadcast chan BoardMessage
}

func (c *BoardController) HandleWebSocket(ctx *gin.Context) {
	boardID := ctx.Param("boardId")
	if boardID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Board ID is required"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}

	c.mutex.Lock()
	hub, exists := c.hubs[boardID]
	if !exists {
		hub = NewBoardHub(boardID, c.service)
		c.hubs[boardID] = hub
		go hub.Run()
	}
	c.mutex.Unlock()

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan BoardMessage, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func NewBoardHub(boardID string, service service.Service) *BoardHub {
	return &BoardHub{
		boardID:   boardID,
		service:   service,
		clients:   make(map[*Client]bool),
		register:  make(chan *Client),
		broadcast: make(chan BoardMessage),
	}
}

func (h *BoardHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.broadcast <- BoardMessage{Type: "USER_DISCONNECTED"}
		c.conn.Close()
	}()

	for {
		var msg BoardMessage
		if err := c.conn.ReadJSON(&msg); err != nil {
			break
		}

		switch msg.Type {
		case "CLIENT_UPDATE":
			board := &domain.Board{
				ID:    c.hub.boardID,
				Nodes: msg.Nodes,
				Edges: msg.Edges,
			}

			if err := c.hub.service.UpdateBoard(context.Background(), board); err != nil {
				continue
			}

			c.hub.broadcast <- BoardMessage{
				Type:  "BOARD_UPDATED",
				Nodes: msg.Nodes,
				Edges: msg.Edges,
			}
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteJSON(message); err != nil {
			break
		}
	}
}

func (c *BoardController) GetBoard(ctx *gin.Context) {
	boardID := ctx.Param("id")
	if boardID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Board ID is required"})
		return
	}

	board, err := c.service.GetBoard(ctx.Request.Context(), boardID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, board)
}

func (c *BoardController) SaveBoard(ctx *gin.Context) {
	var boardData struct {
		Title       string        `json:"title"`
		Description string        `json:"description"`
		Nodes       []domain.Node `json:"nodes"`
		Edges       []domain.Edge `json:"edges"`
	}

	if err := ctx.ShouldBindJSON(&boardData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	board := &domain.Board{
		Title:       boardData.Title,
		Description: boardData.Description,
		Nodes:       boardData.Nodes,
		Edges:       boardData.Edges,
	}

	if err := c.service.SaveBoard(ctx.Request.Context(), board); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "id": board.ID})
}

func (c *BoardController) UpdateBoard(ctx *gin.Context) {
	boardID := ctx.Param("id")
	if boardID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Board ID is required"})
		return
	}

	var boardData struct {
		Title       string        `json:"title"`
		Description string        `json:"description"`
		Nodes       []domain.Node `json:"nodes"`
		Edges       []domain.Edge `json:"edges"`
	}

	if err := ctx.ShouldBindJSON(&boardData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	board := &domain.Board{
		ID:          boardID,
		Title:       boardData.Title,
		Description: boardData.Description,
		Nodes:       boardData.Nodes,
		Edges:       boardData.Edges,
	}

	if err := c.service.UpdateBoard(ctx.Request.Context(), board); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
