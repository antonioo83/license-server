package interfaces

import "github.com/antonioo83/license-server/internal/models"

type LicenseRepository interface {
	MultipleReplace(customerId int, models []models.Licence) error
	DeleteAll(userId int) error
}
