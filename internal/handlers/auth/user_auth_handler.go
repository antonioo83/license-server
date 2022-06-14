package auth

import (
	"fmt"
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"net/http"
	"strings"
)

const Admin = 1
const User = 2

type UserAuth struct {
	Role int
	User models.User
}

type UserAuthHandler struct {
	userRepository interfaces.UserRepository
	config         config.Config
}

func NewUserAuth(userRepository interfaces.UserRepository, config config.Config) *UserAuthHandler {
	return &UserAuthHandler{userRepository, config}
}

func (u UserAuthHandler) GetAuthUser(token string) (*UserAuth, error) {
	if len(token) < 32 {
		return nil, fmt.Errorf("User doesn't have correct a token")
	}

	user, err := u.userRepository.FindByToken(token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("User not found")
	}

	if user.AuthToken == u.config.Auth.AdminAuthToken && token == user.AuthToken {
		return &UserAuth{Role: Admin, User: *user}, nil
	}

	return &UserAuth{Role: User, User: *user}, nil
}

func (u UserAuthHandler) GetToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer", "", 1)

	return strings.ReplaceAll(token, " ", "")
}
