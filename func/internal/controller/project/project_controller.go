package project

import (
	"net/http"

	"func/internal/domain"
	proRepo "func/internal/repository/project"
	"func/internal/service/project"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdProj, err := c.service.CreateProject(proj)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdProj)
}

func (c *ProjectController) Get(ctx *gin.Context) {
	id := ctx.Param("id")

	proj, err := c.service.GetProject(id)
	if err != nil {
		if err == proRepo.ErrProjectNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, proj)
}

func (c *ProjectController) GetByUser(ctx *gin.Context) {
	userID := ctx.Param("userId")

	projects, err := c.service.GetProjectsByUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}

func (c *ProjectController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var proj domain.Project
	if err := ctx.ShouldBindJSON(&proj); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proj.ID = id

	updatedProj, err := c.service.UpdateProject(proj)
	if err != nil {
		if err == proRepo.ErrProjectNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedProj)
}

func (c *ProjectController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeleteProject(id)
	if err != nil {
		if err == proRepo.ErrProjectNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
