package organization

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"func/internal/domain"

	"github.com/google/uuid"
)

type OrganizationRepositoryImpl struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepositoryImpl {
	return &OrganizationRepositoryImpl{db: db}
}

func (r *OrganizationRepositoryImpl) Create(org domain.Organization) (domain.Organization, error) {
	org.ID = uuid.New().String()
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()

	query := `INSERT INTO organizations (id, name, description, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id, name, description, created_at, updated_at`

	err := r.db.QueryRowContext(context.Background(), query,
		org.ID, org.Name, org.Description, org.CreatedAt, org.UpdatedAt).Scan(
		&org.ID, &org.Name, &org.Description, &org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		return domain.Organization{}, err
	}

	return org, nil
}

func (r *OrganizationRepositoryImpl) GetByID(id string) (domain.Organization, error) {
	var org domain.Organization

	query := `SELECT id, name, description, created_at, updated_at 
              FROM organizations WHERE id = $1`

	err := r.db.QueryRowContext(context.Background(), query, id).Scan(
		&org.ID, &org.Name, &org.Description, &org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Organization{}, ErrOrganizationNotFound
		}
		return domain.Organization{}, err
	}

	return org, nil
}

func (r *OrganizationRepositoryImpl) GetAll() ([]domain.Organization, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM organizations`

	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []domain.Organization
	for rows.Next() {
		var org domain.Organization
		if err := rows.Scan(
			&org.ID, &org.Name, &org.Description, &org.CreatedAt, &org.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}

	return orgs, nil
}

func (r *OrganizationRepositoryImpl) Update(org domain.Organization) (domain.Organization, error) {
	org.UpdatedAt = time.Now()

	query := `UPDATE organizations 
              SET name = $1, description = $2, updated_at = $3 
              WHERE id = $4 
              RETURNING id, name, description, created_at, updated_at`

	err := r.db.QueryRowContext(context.Background(), query,
		org.Name, org.Description, org.UpdatedAt, org.ID).Scan(
		&org.ID, &org.Name, &org.Description, &org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Organization{}, ErrOrganizationNotFound
		}
		return domain.Organization{}, err
	}

	return org, nil
}

func (r *OrganizationRepositoryImpl) Delete(id string) error {
	query := `DELETE FROM organizations WHERE id = $1`

	result, err := r.db.ExecContext(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrOrganizationNotFound
	}

	return nil
}
