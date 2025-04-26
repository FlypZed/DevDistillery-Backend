package service

import (
	"errors"
	"fmt"
	"func/internal/domain"
	repository "func/internal/repository/task"
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

func (ts *taskService) GetTasksByProject(projectID string) ([]domain.Task, error) {
	if projectID == "" {
		return nil, errors.New("project ID cannot be empty")
	}

	tasks, err := ts.taskRepository.FindByProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func (ts *taskService) UpdateTaskStatus(taskID string, status string) (*domain.Task, error) {
	if taskID == "" {
		return nil, errors.New("task ID cannot be empty")
	}

	if status == "" {
		return nil, errors.New("status cannot be empty")
	}

	validStatuses := map[string]bool{
		"backlog":     true,
		"todo":        true,
		"in_progress": true,
		"done":        true,
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid task status")
	}

	task, err := ts.taskRepository.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	task.Status = status

	if err := ts.taskRepository.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	return task, nil
}
