package service

import (
	"func/internal/domain"
)

type TaskService interface {
	CreateTask(task *domain.Task) error
	GetTask(id string) (*domain.Task, error)
	UpdateTask(task *domain.Task) error
	DeleteTask(id string) error
}
