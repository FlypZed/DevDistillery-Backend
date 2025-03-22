package repository

import (
	"func/internal/domain"
)

type OrganizationRepository interface {
	Create(organization *domain.Organization) error
	FindByID(id string) (*domain.Organization, error)
	Update(organization *domain.Organization) error
	Delete(id string) error
}
