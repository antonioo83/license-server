package interfaces

import "github.com/antonioo83/license-server/internal/models"

type CustomerRepository interface {
	Replace(userId int, model models.Customer, licenses []models.Licence) error
	Delete(userId int, code string) error
	FindByCode(userId int, code string) (*models.Customer, error)
	IsInDatabase(userId int, code string) (bool, error)
	FindFull(userId int, customerCode string, licenseCode string) (*models.Customer, error)
}
