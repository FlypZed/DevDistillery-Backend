package repository

import (
	"errors"
	"func/internal/domain"
	"gorm.io/gorm"
)

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (tr *teamRepository) Create(team *domain.Team) error {
	if team == nil {
		return errors.New("team is nil")
	}
	return tr.db.Create(team).Error
}

func (tr *teamRepository) FindByID(id string) (*domain.Team, error) {
	var team domain.Team
	if err := tr.db.First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (tr *teamRepository) Update(team *domain.Team) error {
	if team == nil || team.ID == "" {
		return errors.New("team or team ID is nil")
	}
	return tr.db.Save(team).Error
}

func (tr *teamRepository) Delete(id string) error {
	if id == "" {
		return errors.New("team ID is required")
	}
	return tr.db.Delete(&domain.Team{}, "id = ?", id).Error
}
