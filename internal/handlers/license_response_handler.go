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
		activationAt, err := utils.GetTimeFromStr(licenseRequest.ActivationAt)
		if err != nil {
			return nil, err
		}

		expirationAt, err := utils.GetTimeFromStr(licenseRequest.ExpirationAt)
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
		license.CallbackUrl = licenseRequest.CallbackURL
		license.Count = licenseRequest.Count
		license.LicenseKey = licenseKey
		license.ActivationAt = activationAt
		license.ExpirationAt = expirationAt
		license.Duration = int(expirationAt.Sub(activationAt).Hours() / 24)
		licences = append(licences, license)
	}

	return licences, nil
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

func GetDeletedLicenseResponse(r *http.Request, w http.ResponseWriter, param LicenseRouteParameters) {

}

func GetLicenseResponse(r *http.Request, w http.ResponseWriter, param LicenseRouteParameters) {

}

func GetLicensesResponse(r *http.Request, w http.ResponseWriter, param LicenseRouteParameters) {

}
