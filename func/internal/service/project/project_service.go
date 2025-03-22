package service

import (
	"func/internal/domain"
)

type ProjectService interface {
	CreateProject(project *domain.Project) error
	GetProject(id string) (*domain.Project, error)
	UpdateProject(project *domain.Project) error
	DeleteProject(id string) error
}
