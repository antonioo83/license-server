package models

type UserPermission struct {
	ID          int
	UserID      int
	ActionID    int
	ProductType string
	Description string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}
