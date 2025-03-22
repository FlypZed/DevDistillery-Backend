package repository

import (
	"errors"
	"func/internal/domain"
	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (tr *taskRepository) Create(task *domain.Task) error {
	if task == nil {
		return errors.New("task is nil")
	}
	return tr.db.Create(task).Error
}

func (tr *taskRepository) FindByID(id string) (*domain.Task, error) {
	var task domain.Task
	if err := tr.db.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (tr *taskRepository) Update(task *domain.Task) error {
	if task == nil || task.ID == "" {
		return errors.New("task or task ID is nil")
	}
	return tr.db.Save(task).Error
}

func (tr *taskRepository) Delete(id string) error {
	if id == "" {
		return errors.New("task ID is required")
	}
	return tr.db.Delete(&domain.Task{}, "id = ?", id).Error
}
