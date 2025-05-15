package repository

import (
	"func/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id string) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
}
