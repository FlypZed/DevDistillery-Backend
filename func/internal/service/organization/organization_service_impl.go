package organization

import (
	"func/internal/domain"
	"func/internal/repository/organization"
)

type OrganizationServiceImpl struct {
	repo organization.OrganizationRepository
}

func NewOrganizationService(repo organization.OrganizationRepository) *OrganizationServiceImpl {
	return &OrganizationServiceImpl{repo: repo}
}

func (s *OrganizationServiceImpl) CreateOrganization(org domain.Organization) (domain.Organization, error) {
	return s.repo.Create(org)
}

func (s *OrganizationServiceImpl) GetOrganization(id string) (domain.Organization, error) {
	return s.repo.GetByID(id)
}

func (s *OrganizationServiceImpl) GetAllOrganizations() ([]domain.Organization, error) {
	return s.repo.GetAll()
}

func (s *OrganizationServiceImpl) UpdateOrganization(org domain.Organization) (domain.Organization, error) {
	return s.repo.Update(org)
}

func (s *OrganizationServiceImpl) DeleteOrganization(id string) error {
	return s.repo.Delete(id)
}
