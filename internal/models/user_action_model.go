package models

import "time"

type UserAction struct {
	ID          int
	Action      string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
