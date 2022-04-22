package factory

import (
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/handlers/auth"
	repositoryInterfaces "github.com/antonioo83/license-server/internal/repositories/interfaces"
)

func NewUserAuthHandler(userRepository repositoryInterfaces.UserRepository, config config.Config) *auth.UserAuthHandler {
	return auth.NewUserAuth(userRepository, config)
}
