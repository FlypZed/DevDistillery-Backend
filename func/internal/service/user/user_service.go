package service

import (
	"func/internal/domain"
)

type UserService interface {
	CreateUser(user *domain.User) error
	GetUser(id string) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id string) error
}
