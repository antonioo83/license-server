package models

import "time"

type Customer struct {
	ID          int
	UserID      int
	Code        string
	Type        string
	Title       string
	Inn         string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
