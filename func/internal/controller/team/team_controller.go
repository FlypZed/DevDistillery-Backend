package controller

import (
	"func/internal/domain"
	service "func/internal/service/team"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TeamController struct {
	teamService service.TeamService
}

func NewTeamController(teamService service.TeamService) *TeamController {
	return &TeamController{teamService: teamService}
}

func (tc *TeamController) CreateTeam(c *gin.Context) {
	var team domain.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.teamService.CreateTeam(&team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (tc *TeamController) GetTeam(c *gin.Context) {
	id := c.Param("id")
	team, err := tc.teamService.GetTeam(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (tc *TeamController) UpdateTeam(c *gin.Context) {
	id := c.Param("id")
	var team domain.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team.ID = id
	if err := tc.teamService.UpdateTeam(&team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (tc *TeamController) DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	if err := tc.teamService.DeleteTeam(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}
