package domain

import (
	"time"
)

type ProjectStatus string

const (
	ProjectStatusPlanned    ProjectStatus = "planned"
	ProjectStatusInProgress ProjectStatus = "in_progress"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusOnHold     ProjectStatus = "on_hold"
)

type Project struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Status         ProjectStatus `json:"status"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
}
