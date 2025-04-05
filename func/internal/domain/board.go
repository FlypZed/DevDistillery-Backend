package domain

import (
	"time"
)

type Board struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Nodes       []Node    `gorm:"foreignKey:BoardID" json:"nodes"`
	Edges       []Edge    `gorm:"foreignKey:BoardID" json:"edges"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Node struct {
	ID       string         `gorm:"primaryKey" json:"id"`
	BoardID  string         `gorm:"size:36;index" json:"-"`
	Type     string         `json:"type"`
	Position Position       `gorm:"embedded" json:"position"`
	Data     map[string]any `gorm:"serializer:json" json:"data"`
}

type Edge struct {
	ID      string         `gorm:"primaryKey" json:"id"`
	BoardID string         `gorm:"size:36;index" json:"-"`
	Source  string         `json:"source"`
	Target  string         `json:"target"`
	Data    map[string]any `gorm:"serializer:json" json:"data"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
