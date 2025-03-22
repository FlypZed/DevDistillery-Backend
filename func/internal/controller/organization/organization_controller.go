package controller

import (
	"func/internal/domain"
	service "func/internal/service/organization"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrganizationController struct {
	organizationService service.OrganizationService
}

func NewOrganizationController(organizationService service.OrganizationService) *OrganizationController {
	return &OrganizationController{organizationService: organizationService}
}

func (oc *OrganizationController) CreateOrganization(c *gin.Context) {
	var organization domain.Organization
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := oc.organizationService.CreateOrganization(&organization); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, organization)
}

func (oc *OrganizationController) GetOrganization(c *gin.Context) {
	id := c.Param("id")
	organization, err := oc.organizationService.GetOrganization(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, organization)
}

func (oc *OrganizationController) UpdateOrganization(c *gin.Context) {
	id := c.Param("id")
	var organization domain.Organization
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organization.ID = id
	if err := oc.organizationService.UpdateOrganization(&organization); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, organization)
}

func (oc *OrganizationController) DeleteOrganization(c *gin.Context) {
	id := c.Param("id")
	if err := oc.organizationService.DeleteOrganization(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}
