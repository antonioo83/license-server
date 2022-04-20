package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/models"
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
	httpRequest, err := getRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	user.Code = httpRequest.UserId
	user.Role = httpRequest.Role
	user.Title = httpRequest.Title
	user.Description = httpRequest.Description
	userId, err := param.UserRepository.Save(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	actions, err := param.ActionRepository.FindALL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var permissions []models.UserPermission
	for _, product := range httpRequest.Products {
		for _, permission := range product.Permissions {
			action, ok := actions[permission]
			if !ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			var userPermission models.UserPermission
			userPermission.UserID = userId
			userPermission.ActionID = action.ID
			userPermission.ProductType = product.Type
			permissions = append(permissions, userPermission)
		}
	}

	param.PermissionRepository.Replace(permissions)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
}

type Product struct {
	Type        string
	Permissions [4]string
}

type UserRequest struct {
	UserId      string
	Role        string
	Title       string
	Description string
	Products    []Product
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
