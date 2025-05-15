package project

import (
	"log"
	"net/http"
	"strconv"

	"func/internal/domain"
	proRepo "func/internal/repository/project"
	"func/internal/service/project"
	"func/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	service project.ProjectService
}

func NewProjectController(service project.ProjectService) *ProjectController {
	return &ProjectController{service: service}
}

func (c *ProjectController) Create(ctx *gin.Context) {
	var proj domain.Project
	if err := ctx.ShouldBindJSON(&proj); err != nil {
		log.Printf("[ProjectController] Error al bindear JSON: %v\n", err)
		response.Error(ctx, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	token := ctx.GetHeader("Authorization")
	if token == "" {
		token = ctx.Query("token")
	}
	if token == "" {
		response.Error(ctx, http.StatusUnauthorized, "Authorization token required")
		return
	}

	createdProj, err := c.service.CreateProject(proj)
	if err != nil {
		log.Printf("[ProjectController] Error creating project: %v\n", err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to create project: "+err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, createdProj, "Project created successfully")
}

func (c *ProjectController) Get(ctx *gin.Context) {
	id := ctx.Param("id")

	proj, err := c.service.GetProject(id)
	if err != nil {
		if err == proRepo.ErrProjectNotFound {
			response.Error(ctx, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Failed to get project: "+err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, proj, "Project fetched successfully")
}

func (c *ProjectController) GetByUser(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	log.Printf("Received userId param: %s", userIDStr)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Error converting userId to int: %v", err)
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID format")
		return
	}
	log.Printf("Parsed userId: %d", userID)

	projects, err := c.service.GetProjectsByUser(userID)
	if err != nil {
		log.Printf("Error retrieving projects for user %d: %v", userID, err)
		response.Error(ctx, http.StatusInternalServerError, "Failed to get user projects: "+err.Error())
		return
	}

	log.Printf("Found %d projects for user %d", len(projects), userID)
	response.Success(ctx, http.StatusOK, projects, "User projects fetched successfully")
}

func (c *ProjectController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var proj domain.Project
	if err := ctx.ShouldBindJSON(&proj); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}
	proj.ID = id

	updatedProj, err := c.service.UpdateProject(proj)
	if err != nil {
		if err == proRepo.ErrProjectNotFound {
			response.Error(ctx, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Failed to update project: "+err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, updatedProj, "Project updated successfully")
}

func (c *ProjectController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeleteProject(id)
	if err != nil {
		if err == proRepo.ErrProjectNotFound {
			response.Error(ctx, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Failed to delete project: "+err.Error())
		return
	}

	response.Success(ctx, http.StatusNoContent, nil, "Project deleted successfully")
}

func (c *ProjectController) AddMember(ctx *gin.Context) {
	projectID := ctx.Param("id")
	if projectID == "" {
		response.Error(ctx, http.StatusBadRequest, "Project ID is required")
		return
	}

	var req struct {
		UserID string `json:"userId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := c.service.AddMember(projectID, req.UserID); err != nil {
		if err == proRepo.ErrProjectNotFound {
			response.Error(ctx, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Failed to add member: "+err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, nil, "Member added successfully")
}

func (c *ProjectController) RemoveMember(ctx *gin.Context) {
	projectID := ctx.Param("id")
	userID := ctx.Param("userId")

	if projectID == "" || userID == "" {
		response.Error(ctx, http.StatusBadRequest, "Project ID and user ID are required")
		return
	}

	if err := c.service.RemoveMember(projectID, userID); err != nil {
		if err == proRepo.ErrProjectNotFound {
			response.Error(ctx, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Failed to remove member: "+err.Error())
		return
	}

	response.Success(ctx, http.StatusNoContent, nil, "Member removed successfully")
}

func (c *ProjectController) ListMembers(ctx *gin.Context) {
	projectID := ctx.Param("id")
	if projectID == "" {
		response.Error(ctx, http.StatusBadRequest, "Project ID is required")
		return
	}

	members, err := c.service.ListMembers(projectID)
	if err != nil {
		if err == proRepo.ErrProjectNotFound {
			response.Error(ctx, http.StatusNotFound, "Project not found")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Failed to list members: "+err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, members, "Members listed successfully")
}
