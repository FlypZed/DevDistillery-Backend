package repository

import (
	"errors"
	"fmt"
	"func/internal/domain"
	"log"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	return ur.db.Create(user).Error
}

func (ur *userRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	if err := ur.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) FindAll() ([]domain.User, error) {
	var users []domain.User

	if err := ur.db.Table("app_user").Find(&users).Error; err != nil {
		log.Printf("Error en UserRepository.FindAll: %v", err)
		return nil, fmt.Errorf("error al buscar usuarios: %w", err)
	}

	return users, nil
}

func (ur *userRepository) Update(user *domain.User) error {
	if user == nil || user.ID == "" {
		return errors.New("user or user ID is nil")
	}
	return ur.db.Save(user).Error
}

func (ur *userRepository) Delete(id string) error {
	if id == "" {
		return errors.New("user ID is required")
	}
	return ur.db.Delete(&domain.User{}, "id = ?", id).Error
}
