package domain

import "time"

type Entity interface {
	GetID() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

func (u User) GetID() string {
	return u.ID
}

func (u User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}
