package repository

import (
	"errors"
	"func/internal/domain"
	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (pr *projectRepository) Create(project *domain.Project) error {
	if project == nil {
		return errors.New("project is nil")
	}
	return pr.db.Create(project).Error
}

func (pr *projectRepository) FindByID(id string) (*domain.Project, error) {
	var project domain.Project
	if err := pr.db.First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (pr *projectRepository) Update(project *domain.Project) error {
	if project == nil || project.ID == "" {
		return errors.New("project or project ID is nil")
	}
	return pr.db.Save(project).Error
}

func (pr *projectRepository) Delete(id string) error {
	if id == "" {
		return errors.New("project ID is required")
	}
	return pr.db.Delete(&domain.Project{}, "id = ?", id).Error
}
