package project

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"func/internal/domain"

	"github.com/google/uuid"
)

type ProjectRepositoryImpl struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepositoryImpl {
	return &ProjectRepositoryImpl{db: db}
}

func (r *ProjectRepositoryImpl) Create(project domain.Project) (domain.Project, error) {
	project.ID = uuid.New().String()
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	query := `INSERT INTO projects 
              (id, name, description, status, team_id, organization_id, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
              RETURNING id, name, description, status, team_id, organization_id, created_at, updated_at`

	err := r.db.QueryRowContext(context.Background(), query,
		project.ID, project.Name, project.Description, project.Status,
		project.TeamID, project.OrganizationID, project.CreatedAt, project.UpdatedAt).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status,
		&project.TeamID, &project.OrganizationID, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return domain.Project{}, err
	}

	return project, nil
}

func (r *ProjectRepositoryImpl) GetByID(id string) (domain.Project, error) {
	var project domain.Project

	query := `SELECT id, name, description, status, team_id, organization_id, created_at, updated_at 
              FROM projects WHERE id = $1`

	err := r.db.QueryRowContext(context.Background(), query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status,
		&project.TeamID, &project.OrganizationID, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, ErrProjectNotFound
		}
		return domain.Project{}, err
	}

	return project, nil
}

func (r *ProjectRepositoryImpl) GetByUser(userID string) ([]domain.Project, error) {
	query := `SELECT p.id, p.name, p.description, p.status, p.team_id, p.organization_id, p.created_at, p.updated_at
              FROM projects p
              JOIN team_members tm ON p.team_id = tm.team_id
              WHERE tm.user_id = $1`

	rows, err := r.db.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var project domain.Project
		if err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Status,
			&project.TeamID, &project.OrganizationID, &project.CreatedAt, &project.UpdatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *ProjectRepositoryImpl) GetByOrganization(orgID string) ([]domain.Project, error) {
	query := `SELECT id, name, description, status, team_id, organization_id, created_at, updated_at
              FROM projects WHERE organization_id = $1`

	rows, err := r.db.QueryContext(context.Background(), query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var project domain.Project
		if err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Status,
			&project.TeamID, &project.OrganizationID, &project.CreatedAt, &project.UpdatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *ProjectRepositoryImpl) Update(project domain.Project) (domain.Project, error) {
	project.UpdatedAt = time.Now()

	query := `UPDATE projects 
              SET name = $1, description = $2, status = $3, team_id = $4, organization_id = $5, updated_at = $6 
              WHERE id = $7 
              RETURNING id, name, description, status, team_id, organization_id, created_at, updated_at`

	err := r.db.QueryRowContext(context.Background(), query,
		project.Name, project.Description, project.Status, project.TeamID,
		project.OrganizationID, project.UpdatedAt, project.ID).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status,
		&project.TeamID, &project.OrganizationID, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, ErrProjectNotFound
		}
		return domain.Project{}, err
	}

	return project, nil
}

func (r *ProjectRepositoryImpl) AssignTeam(projectID, teamID string) (domain.Project, error) {
	query := `UPDATE projects 
              SET team_id = $1, updated_at = $2 
              WHERE id = $3 
              RETURNING id, name, description, status, team_id, organization_id, created_at, updated_at`

	var project domain.Project
	err := r.db.QueryRowContext(context.Background(), query,
		teamID, time.Now(), projectID).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status,
		&project.TeamID, &project.OrganizationID, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, ErrProjectNotFound
		}
		return domain.Project{}, err
	}

	return project, nil
}

func (r *ProjectRepositoryImpl) Delete(id string) error {
	query := `DELETE FROM projects WHERE id = $1`

	result, err := r.db.ExecContext(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrProjectNotFound
	}

	return nil
}

var ErrProjectNotFound = errors.New("project not found")