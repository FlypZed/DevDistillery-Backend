package service

import (
	"errors"
	"func/internal/domain"
	"func/internal/repository/task"
)

type taskService struct {
	taskRepository repository.TaskRepository
}

func NewTaskService(taskRepository repository.TaskRepository) TaskService {
	return &taskService{taskRepository: taskRepository}
}

func (ts *taskService) CreateTask(task *domain.Task) error {
	if task.Title == "" || task.ProjectID == "" {
		return errors.New("title and project ID are required")
	}

	return ts.taskRepository.Create(task)
}

func (ts *taskService) GetTask(id string) (*domain.Task, error) {
	return ts.taskRepository.FindByID(id)
}

func (ts *taskService) UpdateTask(task *domain.Task) error {
	if task.ID == "" {
		return errors.New("task ID is required")
	}

	return ts.taskRepository.Update(task)
}

func (ts *taskService) DeleteTask(id string) error {
	return ts.taskRepository.Delete(id)
}
