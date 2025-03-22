package repository

import (
	"func/internal/domain"
)

type TeamRepository interface {
	Create(team *domain.Team) error
	FindByID(id string) (*domain.Team, error)
	Update(team *domain.Team) error
	Delete(id string) error
}
