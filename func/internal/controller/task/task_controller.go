package controller

import (
	"func/internal/domain"
	service "func/internal/service/task"
	ws "func/internal/service/websocket"
	"func/pkg/infrastructure"
	"func/pkg/response" // Nuevo paquete importado
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
		response.Error(c, http.StatusBadRequest, "Invalid task data: "+err.Error())
		return
	}

	if err := tc.taskService.CreateTask(&task); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create task: "+err.Error())
		return
	}

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

	response.Success(c, http.StatusOK, task, "Task updated successfully")
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := tc.taskService.DeleteTask(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete task: "+err.Error())
		return
	}

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
		Type: "TASK_UPDATED",
		Data: task,
	})

	response.Success(c, http.StatusOK, task, "Task status updated successfully")
}

func (tc *TaskController) HandleKanbanWebSocket(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.Error(c, http.StatusBadRequest, "Project ID is required")
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
		if token == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization token required")
			return
		}
	} else {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	userID, err := infrastructure.ValidateJWT(token)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
		return
	}

	tc.websocketService.HandleConnection(c.Writer, c.Request, projectID, userID)
}
