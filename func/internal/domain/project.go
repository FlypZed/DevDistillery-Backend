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
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      ProjectStatus `json:"status"`
	Members     []Member      `gorm:"-" json:"members"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

type Member struct {
	UserID    string    `json:"userId" gorm:"primaryKey"`
	ProjectID string    `json:"projectId" gorm:"primaryKey"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joinedAt"`
}
