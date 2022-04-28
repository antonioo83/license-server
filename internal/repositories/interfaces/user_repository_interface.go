package interfaces

import "github.com/antonioo83/license-server/internal/models"

type UserRepository interface {
	Save(user models.User, permissions []models.UserPermission) error
	Update(model models.User, permissions []models.UserPermission) error
	Delete(code string) error
	FindByCode(code string) (*models.User, error)
	FindByToken(code string) (*models.User, error)
	FindALL(limit int, offset int) (*map[int]models.User, error)
	IsInDatabase(code string) (bool, error)
}
