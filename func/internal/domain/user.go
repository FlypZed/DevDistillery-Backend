package domain

import (
	"time"
)

type User struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Picture     string    `json:"picture"`
	GithubID    int64     `json:"githubId"`
	GithubLogin string    `json:"githubLogin"`
	PublicRepos int       `json:"publicRepos"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Password    string    `json:"-" gorm:"-"`
}
