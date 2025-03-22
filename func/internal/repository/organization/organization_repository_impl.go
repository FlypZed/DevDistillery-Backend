package repository

import (
	"errors"
	"func/internal/domain"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (or *organizationRepository) Create(organization *domain.Organization) error {
	if organization == nil {
		return errors.New("organization is nil")
	}
	return or.db.Create(organization).Error
}

func (or *organizationRepository) FindByID(id string) (*domain.Organization, error) {
	var organization domain.Organization
	if err := or.db.First(&organization, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}

func (or *organizationRepository) Update(organization *domain.Organization) error {
	if organization == nil || organization.ID == "" {
		return errors.New("organization or organization ID is nil")
	}
	return or.db.Save(organization).Error
}

func (or *organizationRepository) Delete(id string) error {
	if id == "" {
		return errors.New("organization ID is required")
	}
	return or.db.Delete(&domain.Organization{}, "id = ?", id).Error
}
