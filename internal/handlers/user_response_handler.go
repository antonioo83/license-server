package handlers

import (
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"net/http"
)

type UserRouteParameters struct {
	Config               config.Config
	UserRepository       interfaces.UserRepository
	ActionRepository     interfaces.UserActionRepository
	PermissionRepository interfaces.UserPermissionRepository
}

func GetCreatedUserResponse(r *http.Request, w http.ResponseWriter, param UserRouteParameters) {

}

func GetUpdatedUserResponse(param UserRouteParameters) {

}

func GetDeletedUserResponse(param UserRouteParameters) {

}

func GetUserResponse(param UserRouteParameters) {

}

func GetUsersResponse(param UserRouteParameters) {

}
