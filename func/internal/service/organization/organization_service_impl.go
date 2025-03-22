package service

import (
	"errors"
	"func/internal/domain"
	"func/internal/repository/organization"
)

type organizationService struct {
	organizationRepository repository.OrganizationRepository
}

func NewOrganizationService(organizationRepository repository.OrganizationRepository) OrganizationService {
	return &organizationService{organizationRepository: organizationRepository}
}

func (os *organizationService) CreateOrganization(organization *domain.Organization) error {
	if organization.Name == "" {
		return errors.New("name is required")
	}

	return os.organizationRepository.Create(organization)
}

func (os *organizationService) GetOrganization(id string) (*domain.Organization, error) {
	return os.organizationRepository.FindByID(id)
}

func (os *organizationService) UpdateOrganization(organization *domain.Organization) error {
	if organization.ID == "" {
		return errors.New("organization ID is required")
	}

	return os.organizationRepository.Update(organization)
}

func (os *organizationService) DeleteOrganization(id string) error {
	return os.organizationRepository.Delete(id)
}
