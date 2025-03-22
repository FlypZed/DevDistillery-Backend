package service

import (
	"errors"
	"func/internal/domain"
	"func/internal/repository/team"
)

type teamService struct {
	teamRepository repository.TeamRepository
}

func NewTeamService(teamRepository repository.TeamRepository) TeamService {
	return &teamService{teamRepository: teamRepository}
}

func (ts *teamService) CreateTeam(team *domain.Team) error {
	if team.Name == "" || team.OrganizationID == "" {
		return errors.New("name and organization ID are required")
	}

	return ts.teamRepository.Create(team)
}

func (ts *teamService) GetTeam(id string) (*domain.Team, error) {
	return ts.teamRepository.FindByID(id)
}

func (ts *teamService) UpdateTeam(team *domain.Team) error {
	if team.ID == "" {
		return errors.New("team ID is required")
	}

	return ts.teamRepository.Update(team)
}

func (ts *teamService) DeleteTeam(id string) error {
	return ts.teamRepository.Delete(id)
}
