package project

import (
	"func/internal/domain"
	"func/internal/repository/project"
)

type ProjectServiceImpl struct {
	repo project.ProjectRepository
}

func NewProjectService(repo project.ProjectRepository) *ProjectServiceImpl {
	return &ProjectServiceImpl{repo: repo}
}

func (s *ProjectServiceImpl) CreateProject(project domain.Project) (domain.Project, error) {
	return s.repo.Create(project)
}

func (s *ProjectServiceImpl) GetProject(id string) (domain.Project, error) {
	return s.repo.GetByID(id)
}

func (s *ProjectServiceImpl) GetProjectsByUser(userID string) ([]domain.Project, error) {
	return s.repo.GetByUser(userID)
}

func (s *ProjectServiceImpl) GetProjectsByOrganization(orgID string) ([]domain.Project, error) {
	return s.repo.GetByOrganization(orgID)
}

func (s *ProjectServiceImpl) UpdateProject(project domain.Project) (domain.Project, error) {
	return s.repo.Update(project)
}

func (s *ProjectServiceImpl) AssignTeam(projectID, teamID string) (domain.Project, error) {
	return s.repo.AssignTeam(projectID, teamID)
}

func (s *ProjectServiceImpl) DeleteProject(id string) error {
	return s.repo.Delete(id)
}
