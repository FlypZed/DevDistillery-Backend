package project

import "func/internal/domain"

type ProjectRepository interface {
	Create(project domain.Project) (domain.Project, error)
	GetByID(id string) (domain.Project, error)
	GetByUser(userID string) ([]domain.Project, error)
	Update(project domain.Project) (domain.Project, error)
	Delete(id string) error
}
