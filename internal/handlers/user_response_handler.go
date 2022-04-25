package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/antonioo83/license-server/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"strings"
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

	//userAuth := r.Context().Value("userAuth")

	validate := validator.New()
	err = validate.Struct(httpRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isExist, err := param.UserRepository.IsInDatabase(httpRequest.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isExist {
		http.Error(
			w,
			fmt.Errorf("this user already is exist, orderId=%s", httpRequest.UserId).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	permissions, err := getUserPermissions(httpRequest, param.ActionRepository)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authToken, err := getAuthToken()
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("can't generate user auth token: %w", err).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	var user models.User
	user.Code = httpRequest.UserId
	user.Role = httpRequest.Role
	user.Title = httpRequest.Title
	user.AuthToken = authToken
	user.Description = httpRequest.Description
	err = param.UserRepository.Save(user, permissions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonResponse, err := getJSONResponse("token", authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

func getJSONResponse(key string, value string) ([]byte, error) {
	resp := make(map[string]string)
	resp[key] = value
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return jsonResp, fmt.Errorf("error happened in JSON marshal: %w", err)
	}

	return jsonResp, nil
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

func getAuthToken() (string, error) {
	uuidWithHyphen, err := uuid.NewRandom()
	if err != nil {
		return "", nil
	}
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	return uuid, nil
}

type ProductRequest struct {
	Type        string   `validate:"required,max=50"`
	Permissions []string `validate:"required,oneof='create' 'update' 'delete' 'get'"`
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

func GetUpdatedUserResponse(r *http.Request, w http.ResponseWriter, param UserRouteParameters) {
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

	model, err := param.UserRepository.FindByCode(httpRequest.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if model == nil {
		http.Error(
			w,
			fmt.Errorf("this user isn't exist, userId=%s", httpRequest.UserId).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	permissions, err := getUserPermissions(httpRequest, param.ActionRepository)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authToken := model.AuthToken
	if httpRequest.IsRegenerateToken {
		authToken, err = getAuthToken()
		if err != nil {
			http.Error(
				w,
				fmt.Errorf("can't generate user auth token: %w", err).Error(),
				http.StatusInternalServerError,
			)
			return
		}
	}

	var user models.User
	user.ID = model.ID
	user.Code = httpRequest.UserId
	user.Role = httpRequest.Role
	user.Title = httpRequest.Title
	user.AuthToken = authToken
	user.Description = httpRequest.Description
	err = param.UserRepository.Update(user, permissions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	jsonResponse, err := getJSONResponse("token", authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

type UserDeleteRequest struct {
	UserId string `validate:"required,max=64"`
}

func GetDeletedUserResponse(r *http.Request, w http.ResponseWriter, param UserRouteParameters) {
	httpRequest, err := getDeleteRequest(r)
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

	isExist, err := param.UserRepository.IsInDatabase(httpRequest.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isExist {
		http.Error(
			w,
			fmt.Errorf("this user isn't exist, userId=%s", httpRequest.UserId).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	//userAuth := r.Context().Value("userAuth")
	//u := userAuth.(*auth.UserAuth)
	err = param.UserRepository.Delete(httpRequest.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func getDeleteRequest(r *http.Request) (*UserDeleteRequest, error) {
	var request UserDeleteRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &request, nil
}

func GetUserResponse(param UserRouteParameters) {

}

func GetUsersResponse(param UserRouteParameters) {

}
