package repository

import (
	"func/internal/domain"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	FindByID(id string) (*domain.Task, error)
	Update(task *domain.Task) error
	Delete(id string) error
	FindByProject(projectID string) ([]domain.Task, error)
}
