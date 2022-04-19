package interfaces

import "github.com/antonioo83/license-server/internal/models"

type UserActionRepository interface {
	FindALL() ([]models.UserAction, error)
}
