package service

import (
	"errors"
	"func/internal/domain"
	"func/internal/repository/project"
)

type projectService struct {
	projectRepository repository.ProjectRepository
}

func NewProjectService(projectRepository repository.ProjectRepository) ProjectService {
	return &projectService{projectRepository: projectRepository}
}

func (ps *projectService) CreateProject(project *domain.Project) error {
	if project.Name == "" || project.OrganizationID == "" {
		return errors.New("name and organization ID are required")
	}

	return ps.projectRepository.Create(project)
}

func (ps *projectService) GetProject(id string) (*domain.Project, error) {
	return ps.projectRepository.FindByID(id)
}

func (ps *projectService) UpdateProject(project *domain.Project) error {
	if project.ID == "" {
		return errors.New("project ID is required")
	}

	return ps.projectRepository.Update(project)
}

func (ps *projectService) DeleteProject(id string) error {
	return ps.projectRepository.Delete(id)
}
