package controller

import (
	"func/internal/domain"
	service "func/internal/service/task"
	ws "func/internal/service/websocket"
	"func/pkg/infrastructure"
	"func/pkg/response"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type TaskController struct {
	taskService      service.TaskService
	websocketService ws.Service
}

func NewTaskController(taskService service.TaskService, wsService ws.Service) *TaskController {
	return &TaskController{
		taskService:      taskService,
		websocketService: wsService,
	}
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid task data: "+err.Error())
		return
	}

	if err := tc.taskService.CreateTask(&task); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create task: "+err.Error())
		return
	}

	tc.websocketService.Broadcast(task.ProjectID, ws.Message{
		Type: "TASK_CREATED",
		Data: task,
	})

	response.Success(c, http.StatusCreated, task, "Task created successfully")
}

func (tc *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := tc.taskService.GetTask(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Task not found")
		return
	}

	response.Success(c, http.StatusOK, task, "Task retrieved successfully")
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid task data: "+err.Error())
		return
	}

	task.ID = id
	if err := tc.taskService.UpdateTask(&task); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update task: "+err.Error())
		return
	}

	tc.websocketService.Broadcast(task.ProjectID, ws.Message{
		Type: "TASK_UPDATED",
		Data: task,
	})

	response.Success(c, http.StatusOK, task, "Task updated successfully")
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	task, err := tc.taskService.GetTask(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Task not found")
		return
	}

	if err := tc.taskService.DeleteTask(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete task: "+err.Error())
		return
	}

	tc.websocketService.Broadcast(task.ProjectID, ws.Message{
		Type: "TASK_DELETED",
		Data: map[string]string{"taskId": id},
	})

	response.Success(c, http.StatusOK, nil, "Task deleted successfully")
}

func (tc *TaskController) GetTasksByProject(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.Error(c, http.StatusBadRequest, "Project ID is required")
		return
	}

	tasks, err := tc.taskService.GetTasksByProject(projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get tasks: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, tasks, "Tasks retrieved successfully")
}

func (tc *TaskController) UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.Error(c, http.StatusBadRequest, "Task ID is required")
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid status data: "+err.Error())
		return
	}

	task, err := tc.taskService.UpdateTaskStatus(taskID, statusUpdate.Status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update task status: "+err.Error())
		return
	}

	tc.websocketService.Broadcast(task.ProjectID, ws.Message{
		Type: "TASK_STATUS_UPDATED",
		Data: task,
	})

	response.Success(c, http.StatusOK, task, "Task status updated successfully")
}

func (tc *TaskController) HandleTaskWebSocket(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		log.Printf("[WebSocket] Error: Project ID is required")
		response.Error(c, http.StatusBadRequest, "Project ID is required")
		return
	}

	log.Printf("[WebSocket] Incoming connection for project: %s", projectID)
	log.Printf("[WebSocket] Headers: %v", c.Request.Header)

	if !websocket.IsWebSocketUpgrade(c.Request) {
		log.Printf("[WebSocket] Error: Not a WebSocket upgrade request")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "WebSocket upgrade required"})
		return
	}

	// Get token from Sec-WebSocket-Protocol header
	protocols := websocket.Subprotocols(c.Request)
	log.Printf("[WebSocket] Subprotocols received: %v", protocols)

	var token string
	for _, p := range protocols {
		if strings.HasPrefix(p, "Bearer ") {
			token = strings.TrimPrefix(p, "Bearer ")
			break
		} else if len(p) > 100 {
			token = p
			break
		}
	}

	if token == "" {
		log.Printf("[WebSocket] Error: No token found in subprotocols")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		return
	}

	log.Printf("[WebSocket] Token extracted: %s... (truncated for security)", token[:10])

	userID, err := infrastructure.ValidateJWT(token)
	if err != nil {
		log.Printf("[WebSocket] Error validating JWT: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
		return
	}

	log.Printf("[WebSocket] Authenticated user ID: %s", userID)
	log.Printf("[WebSocket] Proceeding with WebSocket upgrade")

	tc.websocketService.HandleConnection(c.Writer, c.Request, projectID, userID)
}
