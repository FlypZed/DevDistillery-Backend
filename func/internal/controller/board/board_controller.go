package board

import (
	"net/http"
	"strings"

	"func/internal/domain"
	"func/internal/service/board"
	ws "func/internal/service/websocket"
	"func/pkg/infrastructure"

	"github.com/gin-gonic/gin"
)

type BoardController struct {
	service          board.Service
	websocketService ws.Service
}

func NewBoardController(service board.Service, wsService ws.Service) *BoardController {
	return &BoardController{
		service:          service,
		websocketService: wsService,
	}
}

type BoardMessage struct {
	Type  string         `json:"type"`
	Nodes []domain.Node  `json:"nodes,omitempty"`
	Edges []domain.Edge  `json:"edges,omitempty"`
	Data  map[string]any `json:"data,omitempty"`
}

func (c *BoardController) HandleWebSocket(ctx *gin.Context) {
	boardID := ctx.Param("boardId")
	if boardID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Board ID is required"})
		return
	}

	token := ctx.GetHeader("Authorization")
	if token == "" {
		token = ctx.Query("token")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}
	} else {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	userID, err := infrastructure.ValidateJWT(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
		return
	}

	c.websocketService.HandleConnection(ctx.Writer, ctx.Request, boardID, userID)
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

	c.websocketService.Broadcast(board.ID, ws.Message{
		Type: "BOARD_UPDATED",
		Data: board,
	})

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

	c.websocketService.Broadcast(boardID, ws.Message{
		Type: "BOARD_UPDATED",
		Data: board,
	})

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
