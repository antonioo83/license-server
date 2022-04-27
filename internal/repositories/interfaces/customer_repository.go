package interfaces

import "github.com/antonioo83/license-server/internal/models"

type CustomerRepository interface {
	Replace(userId int, model models.Customer, licenses []models.Licence) error
	Delete(code string) error
	FindByCode(code string) (*models.Customer, error)
	IsInDatabase(code string) (bool, error)
}
