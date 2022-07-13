package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/handlers/auth"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/antonioo83/license-server/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type LicenseRouteParameters struct {
	Config             config.Config
	CustomerRepository interfaces.CustomerRepository
	LicenseRepository  interfaces.LicenseRepository
}

type LicenseRequest struct {
	LicenseId    string `validate:"required,max=64"`
	ProductType  string `validate:"required, required,oneof='courier' 'solo' 'pechka54'"`
	CallbackURL  string `validate:"url,max=500"`
	Count        int    `validate:"required,number"`
	LicenseKey   string `validate:"max=500"`
	ActivationAt string `validate:"required,datetime"`
	ExpirationAt string `validate:"required,datetime"`
	Description  string `validate:"max=256"`
}

type CustomerRequest struct {
	CustomerId  string `validate:"required,max=64"`
	Type        string `validate:"required,oneof='service' 'device'"`
	Inn         string `validate:"required,max=12"`
	Title       string `validate:"required,max=100"`
	Description string `validate:"max=256"`
	Licenses    []LicenseRequest
}

func GetReplacedLicenseResponse(r *http.Request, w http.ResponseWriter, param LicenseRouteParameters) {
	httpRequest, err := getCreatedRequest(r)
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

	licenses, err := getCustomerLicenses(httpRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userAuth := r.Context().Value("userAuth").(*auth.UserAuth)

	var customer models.Customer
	customer.UserID = userAuth.User.ID
	customer.Code = httpRequest.CustomerId
	customer.Type = httpRequest.Type
	customer.Title = httpRequest.Title
	customer.Inn = httpRequest.Inn
	customer.Description = httpRequest.Description
	err = param.CustomerRepository.Replace(userAuth.User.ID, customer, licenses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []LicenseResponse
	for _, license := range licenses {
		responses = append(responses, LicenseResponse{LicenseId: license.Code, LicenseKey: license.LicenseKey})
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonResponse, err := getCreateResponse(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

func getCustomerLicenses(httpRequest *CustomerRequest) ([]models.Licence, error) {
	var licences []models.Licence
	for _, licenseRequest := range httpRequest.Licenses {
		activationAt, err := getTimeFromStr(licenseRequest.ActivationAt)
		if err != nil {
			return nil, err
		}

		expirationAt, err := getTimeFromStr(licenseRequest.ExpirationAt)
		if err != nil {
			return nil, err
		}

		licenseKey := licenseRequest.LicenseKey
		if licenseKey == "" {
			licenseKey, err = getLicenseKey()
			if err != nil {
				return nil, err
			}
		}

		var license models.Licence
		license.Code = licenseRequest.LicenseId
		license.ProductType = licenseRequest.ProductType
		license.CallbackURL = licenseRequest.CallbackURL
		license.Count = licenseRequest.Count
		license.LicenseKey = licenseKey
		license.ActivationAt = activationAt
		license.ExpirationAt = expirationAt
		license.Duration = int(expirationAt.Sub(activationAt).Hours() / 24)
		license.Description = licenseRequest.Description
		licences = append(licences, license)
	}

	return licences, nil
}

func getTimeFromStr(dateTime string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	time, err := time.Parse(layout, dateTime)
	if err != nil {
		return time, err
	}

	return time, nil
}

func getLicenseKey() (string, error) {
	uuidWithHyphen, err := uuid.NewRandom()
	if err != nil {
		return "", nil
	}
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	return uuid, nil
}

func getCreatedRequest(r *http.Request) (*CustomerRequest, error) {
	var request CustomerRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &request, nil
}

type LicenseResponse struct {
	LicenseId  string
	LicenseKey string
}

func getCreateResponse(responses []LicenseResponse) ([]byte, error) {
	jsonResp, err := json.Marshal(responses)
	if err != nil {
		return jsonResp, fmt.Errorf("error happened in JSON marshal: %w", err)
	}

	return jsonResp, nil
}

type LicenseDeleteRequest struct {
	CustomerId string `validate:"required,max=64"`
	LicenseId  string `validate:"required,max=64"`
}

func GetDeletedLicenseResponse(r *http.Request, w http.ResponseWriter, param LicenseRouteParameters) {
	httpRequest, err := getDeleteLicenseRequest(r)
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

	userAuth := r.Context().Value("userAuth").(*auth.UserAuth)

	customer, err := param.CustomerRepository.FindByCode(userAuth.User.ID, httpRequest.CustomerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if customer == nil {
		http.Error(
			w,
			fmt.Errorf("this customer isn't exist, customerId=%s", httpRequest.CustomerId).Error(),
			http.StatusInternalServerError,
		)
		return
	}

	err = param.LicenseRepository.Delete(customer.ID, httpRequest.LicenseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func getDeleteLicenseRequest(r *http.Request) (*LicenseDeleteRequest, error) {
	var request LicenseDeleteRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &request, nil
}

type CustomerGetRequest struct {
	CustomerId string `validate:"required,min=1,max=64"`
	LicenseId  string `validate:"max=64"`
}

func GetLicenseResponse(r *http.Request, w http.ResponseWriter, param LicenseRouteParameters) {
	httpRequest := getLicenseRequest(r)
	validate := validator.New()
	err := validate.Struct(httpRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value("userAuth").(*auth.UserAuth)
	customer, err := param.CustomerRepository.FindFull(userAuth.User.ID, httpRequest.CustomerId, httpRequest.LicenseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := getLicenseJsonResponse(getCustomerResponse(*customer))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.LogErr(w.Write(jsonResponse))
}

func getLicenseRequest(r *http.Request) *CustomerGetRequest {
	var request CustomerGetRequest
	request.CustomerId = r.URL.Query().Get("customerId")
	request.LicenseId = r.URL.Query().Get("licenseId")

	return &request
}

type LicenseGetResponse struct {
	LicenseId    string `json:"licenseId"`
	ProductType  string `json:"productType"`
	CallbackURL  string `json:"CallbackURL"`
	Count        int    `json:"count"`
	LicenseKey   string `json:"licenseKey"`
	ActivationAt string `json:"activationAt"`
	ExpirationAt string `json:"expirationAt"`
	Description  string `json:"description"`
}

type CustomerGetResponse struct {
	CustomerId  string               `json:"customerId"`
	Type        string               `json:"type"`
	Inn         string               `json:"inn"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Licenses    []LicenseGetResponse `json:"licenses"`
}

func getCustomerResponse(user models.Customer) CustomerGetResponse {
	var customers = make(map[int]models.Customer)
	customers[user.ID] = user
	responses := getCustomerResponses(&customers)
	for _, response := range responses {
		return response
	}

	return CustomerGetResponse{}
}

func getCustomerResponses(customers *map[int]models.Customer) []CustomerGetResponse {
	var responses []CustomerGetResponse
	for _, customer := range *customers {
		var response CustomerGetResponse
		response.CustomerId = customer.Code
		response.Type = customer.Type
		response.Inn = customer.Inn
		response.Title = customer.Title
		response.Description = customer.Description
		var licenseResponse LicenseGetResponse
		for _, license := range customer.Licenses {
			licenseResponse.LicenseId = license.Code
			licenseResponse.ProductType = license.ProductType
			licenseResponse.CallbackURL = license.CallbackURL
			licenseResponse.Count = license.Count
			licenseResponse.LicenseKey = license.LicenseKey
			licenseResponse.ActivationAt = license.ActivationAt.Format("2006-01-02 15:04:05")
			licenseResponse.ExpirationAt = license.ExpirationAt.Format("2006-01-02 15:04:05")
			licenseResponse.Description = license.Description
			response.Licenses = append(response.Licenses, licenseResponse)
		}
		responses = append(responses, response)
	}

	return responses
}

func getLicenseJsonResponse(resp CustomerGetResponse) ([]byte, error) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return jsonResp, fmt.Errorf("error happened in JSON marshal: %w", err)
	}

	return jsonResp, nil
}
