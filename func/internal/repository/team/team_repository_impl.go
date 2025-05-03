package team

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"func/internal/domain"

	"github.com/google/uuid"
)

type TeamRepositoryImpl struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepositoryImpl {
	return &TeamRepositoryImpl{db: db}
}

func (r *TeamRepositoryImpl) Create(team domain.Team) (domain.Team, error) {
	team.ID = uuid.New().String()
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()

	query := `INSERT INTO teams (id, name, description, organization_id, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              RETURNING id, name, description, organization_id, created_at, updated_at`

	err := r.db.QueryRowContext(context.Background(), query,
		team.ID, team.Name, team.Description, team.OrganizationID, team.CreatedAt, team.UpdatedAt).Scan(
		&team.ID, &team.Name, &team.Description, &team.OrganizationID, &team.CreatedAt, &team.UpdatedAt)

	if err != nil {
		return domain.Team{}, err
	}

	return team, nil
}

func (r *TeamRepositoryImpl) GetByID(id string) (domain.Team, error) {
	var team domain.Team

	query := `SELECT id, name, description, organization_id, created_at, updated_at 
              FROM teams WHERE id = $1`

	err := r.db.QueryRowContext(context.Background(), query, id).Scan(
		&team.ID, &team.Name, &team.Description, &team.OrganizationID, &team.CreatedAt, &team.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Team{}, ErrTeamNotFound
		}
		return domain.Team{}, err
	}

	return team, nil
}

func (r *TeamRepositoryImpl) GetByOrganization(orgID string) ([]domain.Team, error) {
	query := `SELECT id, name, description, organization_id, created_at, updated_at 
              FROM teams WHERE organization_id = $1`

	rows, err := r.db.QueryContext(context.Background(), query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var team domain.Team
		if err := rows.Scan(
			&team.ID, &team.Name, &team.Description, &team.OrganizationID, &team.CreatedAt, &team.UpdatedAt,
		); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (r *TeamRepositoryImpl) Update(team domain.Team) (domain.Team, error) {
	team.UpdatedAt = time.Now()

	query := `UPDATE teams 
              SET name = $1, description = $2, organization_id = $3, updated_at = $4 
              WHERE id = $5 
              RETURNING id, name, description, organization_id, created_at, updated_at`

	err := r.db.QueryRowContext(context.Background(), query,
		team.Name, team.Description, team.OrganizationID, team.UpdatedAt, team.ID).Scan(
		&team.ID, &team.Name, &team.Description, &team.OrganizationID, &team.CreatedAt, &team.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Team{}, ErrTeamNotFound
		}
		return domain.Team{}, err
	}

	return team, nil
}

func (r *TeamRepositoryImpl) Delete(id string) error {
	query := `DELETE FROM teams WHERE id = $1`

	result, err := r.db.ExecContext(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTeamNotFound
	}

	return nil
}

var ErrTeamNotFound = errors.New("team not found")
