package interfaces

import "github.com/antonioo83/license-server/internal/models"

type LicenseRepository interface {
	MultipleReplace(customerId int, models []models.Licence) error
	DeleteAll(customerId int) error
	Delete(customerId int, code string) error
	FindByCode(code string) (*models.Licence, error)
	FindAllExpired(maxAttempts uint, limit uint, offset uint) ([]models.Licence, error)
	UpdateCallbackOptions(licenseId int, isSentCallback uint, callbackAttempts uint) error
	IsInDatabase(code string) (bool, error)
}
