package team

import (
	"net/http"

	"func/internal/domain"
	teamRepo "func/internal/repository/team"
	"func/internal/service/team"

	"github.com/gin-gonic/gin"
)

type TeamController struct {
	service team.TeamService
}

func NewTeamController(service team.TeamService) *TeamController {
	return &TeamController{service: service}
}

func (c *TeamController) Create(ctx *gin.Context) {
	var team domain.Team
	if err := ctx.ShouldBindJSON(&team); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTeam, err := c.service.CreateTeam(team)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdTeam)
}

func (c *TeamController) Get(ctx *gin.Context) {
	id := ctx.Param("id")

	team, err := c.service.GetTeam(id)
	if err != nil {
		if err == teamRepo.ErrTeamNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, team)
}

func (c *TeamController) GetByOrganization(ctx *gin.Context) {
	orgID := ctx.Param("orgId")

	teams, err := c.service.GetTeamsByOrganization(orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, teams)
}

func (c *TeamController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var team domain.Team
	if err := ctx.ShouldBindJSON(&team); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	team.ID = id

	updatedTeam, err := c.service.UpdateTeam(team)
	if err != nil {
		if err == teamRepo.ErrTeamNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedTeam)
}

func (c *TeamController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeleteTeam(id)
	if err != nil {
		if err == teamRepo.ErrTeamNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
