package models

import "time"

type User struct {
	ID          int
	Code        string
	Role        string
	Title       string
	AuthToken   string
	Action      string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
