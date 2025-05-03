package team

import (
	"func/internal/domain"
	"func/internal/repository/team"
)

type TeamServiceImpl struct {
	repo team.TeamRepository
}

func NewTeamService(repo team.TeamRepository) *TeamServiceImpl {
	return &TeamServiceImpl{repo: repo}
}

func (s *TeamServiceImpl) CreateTeam(team domain.Team) (domain.Team, error) {
	return s.repo.Create(team)
}

func (s *TeamServiceImpl) GetTeam(id string) (domain.Team, error) {
	return s.repo.GetByID(id)
}

func (s *TeamServiceImpl) GetTeamsByOrganization(orgID string) ([]domain.Team, error) {
	return s.repo.GetByOrganization(orgID)
}

func (s *TeamServiceImpl) UpdateTeam(team domain.Team) (domain.Team, error) {
	return s.repo.Update(team)
}

func (s *TeamServiceImpl) DeleteTeam(id string) error {
	return s.repo.Delete(id)
}
