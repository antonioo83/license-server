package interfaces

import "github.com/antonioo83/license-server/internal/models"

type UserPermissionRepository interface {
	MultipleInsert(userId int, models []models.UserPermission) error
	FindALL(userId int) ([]models.UserPermission, error)
	Delete(userId int) error
}
