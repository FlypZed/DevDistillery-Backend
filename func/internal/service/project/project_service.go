package project

import "func/internal/domain"

type ProjectService interface {
	CreateProject(project domain.Project) (domain.Project, error)
	GetProject(id string) (domain.Project, error)
	GetProjectsByUser(userID int) ([]domain.Project, error)
	UpdateProject(project domain.Project) (domain.Project, error)
	DeleteProject(id string) error
	AddMember(projectID, userID string) error
	RemoveMember(projectID, userID string) error
	ListMembers(projectID string) ([]domain.Member, error)
}
