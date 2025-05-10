package project

import (
	"errors"
	"fmt"
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
	createdProject, err := s.repo.Create(project)
	if err != nil {
		return domain.Project{}, err
	}

	if project.CreatedBy != "" {
		if err := s.repo.AddMember(createdProject.ID, project.CreatedBy); err != nil {
			_ = s.repo.Delete(createdProject.ID)
			return domain.Project{}, fmt.Errorf("failed to add creator as member: %v", err)
		}

		members, err := s.repo.ListMembers(createdProject.ID)
		if err == nil {
			createdProject.Members = members
		}
	}

	return createdProject, nil
}

func (s *ProjectServiceImpl) GetProject(id string) (domain.Project, error) {
	return s.repo.GetByID(id)
}

func (s *ProjectServiceImpl) GetProjectsByUser(userID string) ([]domain.Project, error) {
	return s.repo.GetByUser(userID)
}

func (s *ProjectServiceImpl) UpdateProject(project domain.Project) (domain.Project, error) {
	return s.repo.Update(project)
}

func (s *ProjectServiceImpl) DeleteProject(id string) error {
	return s.repo.Delete(id)
}

func (s *ProjectServiceImpl) AddMember(projectID, userID string) error {
	if projectID == "" || userID == "" {
		return errors.New("project ID and user ID are required")
	}
	return s.repo.AddMember(projectID, userID)
}

func (s *ProjectServiceImpl) RemoveMember(projectID, userID string) error {
	if projectID == "" || userID == "" {
		return errors.New("project ID and user ID are required")
	}
	return s.repo.RemoveMember(projectID, userID)
}

func (s *ProjectServiceImpl) ListMembers(projectID string) ([]domain.Member, error) {
	if projectID == "" {
		return nil, errors.New("project ID is required")
	}
	return s.repo.ListMembers(projectID)
}
