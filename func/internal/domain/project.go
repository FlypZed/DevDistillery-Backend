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
	CreatedBy   string        `json:"createdBy"`
	Members     []Member      `gorm:"-" json:"members"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

type Member struct {
	UserID    string    `json:"userId" gorm:"-"`
	ProjectID string    `json:"projectId" gorm:"-"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joinedAt"`
	Name      string    `json:"name" gorm:"-"`
	Picture   string    `json:"picture" gorm:"-"`
}
