package organization

import (
	"net/http"

	"func/internal/domain"
	orgRepo "func/internal/repository/organization"
	"func/internal/service/organization"

	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	service organization.OrganizationService
}

func NewOrganizationController(service organization.OrganizationService) *OrganizationController {
	return &OrganizationController{service: service}
}

func (c *OrganizationController) Create(ctx *gin.Context) {
	var org domain.Organization
	if err := ctx.ShouldBindJSON(&org); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if org.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	createdOrg, err := c.service.CreateOrganization(org)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdOrg)
}

func (c *OrganizationController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	org, err := c.service.GetOrganization(id)
	if err != nil {
		if err == orgRepo.ErrOrganizationNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, org)
}

func (c *OrganizationController) GetAll(ctx *gin.Context) {
	orgs, err := c.service.GetAllOrganizations()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch organizations"})
		return
	}

	if len(orgs) == 0 {
		ctx.JSON(http.StatusOK, []domain.Organization{})
		return
	}

	ctx.JSON(http.StatusOK, orgs)
}

func (c *OrganizationController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var org domain.Organization
	if err := ctx.ShouldBindJSON(&org); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	org.ID = id
	updatedOrg, err := c.service.UpdateOrganization(org)
	if err != nil {
		if err == orgRepo.ErrOrganizationNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedOrg)
}

func (c *OrganizationController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	err := c.service.DeleteOrganization(id)
	if err != nil {
		if err == orgRepo.ErrOrganizationNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
