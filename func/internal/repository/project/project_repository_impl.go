package project

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	query := `INSERT INTO project 
              (id, name, description, status, created_by, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) 
              RETURNING id, name, description, status, created_by, created_at, updated_at`

	err := r.db.QueryRowContext(context.Background(), query,
		project.ID, project.Name, project.Description, project.Status, project.CreatedBy,
		project.CreatedAt, project.UpdatedAt).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status, &project.CreatedBy,
		&project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return domain.Project{}, err
	}

	return project, nil
}

func (r *ProjectRepositoryImpl) GetByID(id string) (domain.Project, error) {
	var project domain.Project

	query := `SELECT id, name, description, status, created_at, updated_at 
              FROM project WHERE id = $1`

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

func (r *ProjectRepositoryImpl) GetByUser(userID int) ([]domain.Project, error) {
	query := `SELECT p.id, p.name, p.description, p.status, p.created_by, p.created_at, p.updated_at
              FROM project p
              LEFT JOIN project_member pm ON p.id = pm.project_id
              WHERE pm.user_id = $1 OR p.created_by = $1
              GROUP BY p.id`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedBy,
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

	query := `UPDATE project
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
	query := `DELETE FROM project WHERE id = $1`

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

func (r *ProjectRepositoryImpl) AddMember(projectID, userID string) error {
	query := `INSERT INTO project_member (project_id, user_id, role, joined_at) 
              VALUES ($1, $2, 'member', NOW()) 
              ON CONFLICT (project_id, user_id) DO NOTHING`
	_, err := r.db.Exec(query, projectID, userID)
	return err
}

func (r *ProjectRepositoryImpl) RemoveMember(projectID, userID string) error {
	query := `DELETE FROM project_member WHERE project_id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, projectID, userID)
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

func (r *ProjectRepositoryImpl) ListMembers(projectID string) ([]domain.Member, error) {
	query := `
        SELECT 
            pm.user_id::text, 
            pm.project_id::text, 
            pm.role, 
            pm.joined_at, 
            u.name, 
            u.picture
        FROM project_member pm
        JOIN app_user u ON pm.user_id = u.id
        WHERE pm.project_id = $1::uuid`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var members []domain.Member
	for rows.Next() {
		var m domain.Member
		err := rows.Scan(&m.UserID, &m.ProjectID, &m.Role, &m.JoinedAt, &m.Name, &m.Picture)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		members = append(members, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %v", err)
	}

	return members, nil
}

var ErrProjectNotFound = errors.New("project not found")
