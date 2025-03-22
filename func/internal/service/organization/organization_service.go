package service

import (
	"func/internal/domain"
)

type OrganizationService interface {
	CreateOrganization(organization *domain.Organization) error
	GetOrganization(id string) (*domain.Organization, error)
	UpdateOrganization(organization *domain.Organization) error
	DeleteOrganization(id string) error
}
