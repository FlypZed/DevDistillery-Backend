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
              (id, name, description, status, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              RETURNING id, name, description, status, created_at, updated_at`

	err := r.db.QueryRowContext(context.Background(), query,
		project.ID, project.Name, project.Description, project.Status,
		project.CreatedAt, project.UpdatedAt).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status,
		&project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return domain.Project{}, err
	}

	return project, nil
}

func (r *ProjectRepositoryImpl) GetByID(id string) (domain.Project, error) {
	var project domain.Project

	query := `SELECT id, name, description, status, created_at, updated_at 
              FROM projects WHERE id = $1`

	err := r.db.QueryRowContext(context.Background(), query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status,
		&project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, ErrProjectNotFound
		}
		return domain.Project{}, err
	}

	return project, nil
}

func (r *ProjectRepositoryImpl) GetByUser(userID string) ([]domain.Project, error) {
	query := `SELECT p.id, p.name, p.description, p.status, p.created_at, p.updated_at
              FROM projects p
              JOIN user_projects up ON p.id = up.project_id
              WHERE up.user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Status,
			&p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepositoryImpl) Update(project domain.Project) (domain.Project, error) {
	project.UpdatedAt = time.Now()

	query := `UPDATE projects 
              SET name = $1, description = $2, status = $3, updated_at = $4 
              WHERE id = $5 
              RETURNING id, name, description, status, created_at, updated_at`

	err := r.db.QueryRowContext(context.Background(), query,
		project.Name, project.Description, project.Status, project.UpdatedAt, project.ID).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status,
		&project.CreatedAt, &project.UpdatedAt)

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
