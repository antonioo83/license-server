package interfaces

import "github.com/antonioo83/license-server/internal/models"

type UserRepository interface {
	Save(model models.User) (int, error)
	Update(model models.User) error
	Delete(code string) error
	FindByCode(code string) (*models.User, error)
	FindALL(limit int, offset int) (*map[string]models.User, error)
	IsInDatabase(code string) (bool, error)
}
