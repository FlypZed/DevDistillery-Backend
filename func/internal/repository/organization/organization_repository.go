package organization

import (
	"errors"
	"func/internal/domain"
)

type OrganizationRepository interface {
	Create(org domain.Organization) (domain.Organization, error)
	GetByID(id string) (domain.Organization, error)
	GetAll() ([]domain.Organization, error)
	Update(org domain.Organization) (domain.Organization, error)
	Delete(id string) error
}

var ErrOrganizationNotFound = errors.New("organization not found")