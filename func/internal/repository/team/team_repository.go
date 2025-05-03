package team

import "func/internal/domain"

type TeamRepository interface {
	Create(team domain.Team) (domain.Team, error)
	GetByID(id string) (domain.Team, error)
	GetByOrganization(orgID string) ([]domain.Team, error)
	Update(team domain.Team) (domain.Team, error)
	Delete(id string) error
}
