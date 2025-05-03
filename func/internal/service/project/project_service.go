package project

import "func/internal/domain"

type ProjectService interface {
    CreateProject(project domain.Project) (domain.Project, error)
    GetProject(id string) (domain.Project, error)
    GetProjectsByUser(userID string) ([]domain.Project, error)
    GetProjectsByOrganization(orgID string) ([]domain.Project, error)
    UpdateProject(project domain.Project) (domain.Project, error)
    AssignTeam(projectID, teamID string) (domain.Project, error)
    DeleteProject(id string) error
}