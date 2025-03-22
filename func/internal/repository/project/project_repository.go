package repository

import (
	"func/internal/domain"
)

type ProjectRepository interface {
	Create(project *domain.Project) error
	FindByID(id string) (*domain.Project, error)
	Update(project *domain.Project) error
	Delete(id string) error
}
