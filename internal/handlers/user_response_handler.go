package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type UserRouteParameters struct {
	Config               config.Config
	UserRepository       interfaces.UserRepository
	ActionRepository     interfaces.UserActionRepository
	PermissionRepository interfaces.UserPermissionRepository
}

func GetCreatedUserResponse(r *http.Request, w http.ResponseWriter, param UserRouteParameters) {
	httpRequest, err := getRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(httpRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	permissions, err := getUserPermissions(httpRequest, param.ActionRepository)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user models.User
	user.Code = httpRequest.UserId
	user.Role = httpRequest.Role
	user.Title = httpRequest.Title
	user.Description = httpRequest.Description
	err = param.UserRepository.Save(user, permissions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
}

func getUserPermissions(httpRequest *UserRequest, aRep interfaces.UserActionRepository) ([]models.UserPermission, error) {
	actions, err := aRep.FindALL()
	if err != nil {
		return nil, fmt.Errorf("can't load actions from db: %w", err)
	}

	var permissions []models.UserPermission
	for _, product := range httpRequest.Products {
		for _, permission := range product.Permissions {
			action, ok := actions[permission]
			if !ok {
				return nil, fmt.Errorf("can't find permission: %s", permission)
			}

			var userPermission models.UserPermission
			userPermission.ActionID = action.ID
			userPermission.ProductType = product.Type
			permissions = append(permissions, userPermission)
		}
	}

	return permissions, nil
}

type ProductRequest struct {
	Type        string    `validate:"required,max=50"`
	Permissions [4]string `validate:"required,oneof='create' 'update' 'delete' 'get'"`
}

type UserRequest struct {
	UserId            string `validate:"required,max=64"`
	Role              string `validate:"required,oneof='service' 'device'"`
	Title             string `validate:"required,max=100"`
	Description       string `validate:"max=256"`
	Products          []ProductRequest
	IsRegenerateToken bool
}

func getRequest(r *http.Request) (*UserRequest, error) {
	var request UserRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &request, nil
}

func GetUpdatedUserResponse(param UserRouteParameters) {

}

func GetDeletedUserResponse(param UserRouteParameters) {

}

func GetUserResponse(param UserRouteParameters) {

}

func GetUsersResponse(param UserRouteParameters) {

}
