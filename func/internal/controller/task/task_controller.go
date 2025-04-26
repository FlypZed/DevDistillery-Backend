package controller

import (
	"func/internal/domain"
	service "func/internal/service/task"
	ws "func/internal/service/websocket"
	"func/pkg/infrastructure"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.taskService.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (tc *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := tc.taskService.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = id
	if err := tc.taskService.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := tc.taskService.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (tc *TaskController) GetTasksByProject(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	tasks, err := tc.taskService.GetTasksByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := tc.taskService.UpdateTaskStatus(taskID, statusUpdate.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tc.websocketService.Broadcast(task.ProjectID, ws.Message{
		Type: "TASK_UPDATED",
		Data: task,
	})

	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) HandleKanbanWebSocket(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}
	} else {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	userID, err := infrastructure.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
		return
	}

	tc.websocketService.HandleConnection(c.Writer, c.Request, projectID, userID)
}
