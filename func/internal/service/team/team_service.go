package service

import (
	"func/internal/domain"
)

type TeamService interface {
	CreateTeam(team *domain.Team) error
	GetTeam(id string) (*domain.Team, error)
	UpdateTeam(team *domain.Team) error
	DeleteTeam(id string) error
}
