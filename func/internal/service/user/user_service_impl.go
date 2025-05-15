package service

import (
	"errors"
	"func/internal/domain"
	repository "func/internal/repository/user"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (us *userService) CreateUser(user *domain.User) error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("name, email, and password are required")
	}

	return us.userRepository.Create(user)
}

func (us *userService) GetUser(id string) (*domain.User, error) {
	return us.userRepository.FindByID(id)
}

func (us *userService) GetAllUsers() ([]domain.User, error) {
	return us.userRepository.FindAll()
}

func (us *userService) UpdateUser(user *domain.User) error {
	if user.ID == "" {
		return errors.New("user ID is required")
	}

	return us.userRepository.Update(user)
}

func (us *userService) DeleteUser(id string) error {
	return us.userRepository.Delete(id)
}
