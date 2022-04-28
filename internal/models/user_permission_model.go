package models

import "time"

type UserPermission struct {
	ID          int
	UserID      int
	ActionID    int
	ProductType string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Action      UserAction
}
