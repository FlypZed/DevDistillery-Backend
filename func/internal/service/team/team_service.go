package team

import "func/internal/domain"

type TeamService interface {
	CreateTeam(team domain.Team) (domain.Team, error)
	GetTeam(id string) (domain.Team, error)
	GetTeamsByOrganization(orgID string) ([]domain.Team, error)
	UpdateTeam(team domain.Team) (domain.Team, error)
	DeleteTeam(id string) error
}
