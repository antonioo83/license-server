package services

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	callbackInterfaces "github.com/antonioo83/license-server/internal/services/interfaces"
	"github.com/antonioo83/license-server/internal/utils"
	"github.com/jinzhu/copier"
	"log"
	"net/http"
	"strings"
	"time"
)

type licenseCallbackService struct {
	config config.Callback
	rep    interfaces.LicenseRepository
}

func NewLicenseCallbackService(config config.Callback, rep interfaces.LicenseRepository) callbackInterfaces.LicenseCallbackService {
	return &licenseCallbackService{config, rep}
}

func (c licenseCallbackService) SendCallbacks() {
	licenses, err := c.rep.FindAllExpired(c.config.MaxAttempts, c.config.LimitUnitOfTime, 0)
	if err != nil {
		log.Printf("I can't get licenses for send callbacks: %v", err)
		return
	}
	callbackRequests, err := c.getLicenseCallbackRequests(licenses)
	if err != nil {
		log.Printf("I can't convert license models to callback requests: %v", err)
		return
	}
	c.runSendCallbackWorker(callbackRequests, c.rep)
}

type LicenseCallbackRequest struct {
	CustomerCode     string    `json:"customerId" copier:"Customer.Code"`
	LicenseCode      string    `json:"licenseId" copier:"Code"`
	ProductType      string    `json:"productType"`
	CallbackURL      string    `json:"callbackUrl"`
	Count            int       `json:"count"`
	LicenseKey       string    `json:"licenseKey"`
	ActivationAt     time.Time `json:"activationAt"`
	ExpirationAt     time.Time `json:"expirationAt"`
	Description      string    `json:"description"`
	LicenseId        int       `json:"-" copier:"ID"`
	CallbackAttempts uint      `json:"-"`
	UserAuthToken    string    `json:"-" copier:"Customer.User.AuthToken"`
}

func (c licenseCallbackService) getLicenseCallbackRequests(licences []models.Licence) ([]LicenseCallbackRequest, error) {
	var sendCallbacks []LicenseCallbackRequest
	for _, licence := range licences {
		var licenseCallbackRequest LicenseCallbackRequest
		err := copier.Copy(&licenseCallbackRequest, &licence)
		if err != nil {
			return nil, err
		}
		licenseCallbackRequest.LicenseId = licence.ID
		licenseCallbackRequest.CustomerCode = licence.Customer.Code
		licenseCallbackRequest.LicenseCode = licence.Code
		licenseCallbackRequest.UserAuthToken = licence.Customer.User.AuthToken
		sendCallbacks = append(sendCallbacks, licenseCallbackRequest)
	}

	return sendCallbacks, nil
}

func (c licenseCallbackService) runSendCallbackWorker(sendCallbacks []LicenseCallbackRequest, licenseRep interfaces.LicenseRepository) {
	for _, sendCallback := range sendCallbacks {
		go func(sendCallback LicenseCallbackRequest, licenseRep interfaces.LicenseRepository) {
			httpStatus, err := c.sendRequest(sendCallback)
			if err != nil {
				fmt.Printf("i can't send callback request: %v\n", err)
				return
			}
			if httpStatus == http.StatusCreated {
				err = licenseRep.UpdateCallbackOptions(sendCallback.LicenseId, 1, 1)
				if err != nil {
					fmt.Printf("i can't send callback request: %v\n", err)
				}
				return
			} else {
				err = licenseRep.UpdateCallbackOptions(sendCallback.LicenseId, 0, sendCallback.CallbackAttempts+1)
				if err != nil {
					fmt.Printf("i can't send callback request: %v\n", err)
				}
				return
			}
		}(sendCallback, licenseRep)
	}

	return
}

func (c licenseCallbackService) sendRequest(callback LicenseCallbackRequest) (int, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	jsonRequest, err := c.getJSONRequest(callback)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("POST", callback.CallbackURL, strings.NewReader(string(jsonRequest)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", callback.UserAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer utils.ResourceClose(resp.Body)

	return resp.StatusCode, nil
}

func (c licenseCallbackService) getJSONRequest(request LicenseCallbackRequest) ([]byte, error) {
	jsonResp, err := json.Marshal(request)
	if err != nil {
		return []byte(""), fmt.Errorf("i can't decode json request: %w", err)
	}

	return jsonResp, nil
}
