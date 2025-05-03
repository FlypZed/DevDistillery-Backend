package organization

import "func/internal/domain"

type OrganizationService interface {
	CreateOrganization(org domain.Organization) (domain.Organization, error)
	GetOrganization(id string) (domain.Organization, error)
	GetAllOrganizations() ([]domain.Organization, error)
	UpdateOrganization(org domain.Organization) (domain.Organization, error)
	DeleteOrganization(id string) error
}
