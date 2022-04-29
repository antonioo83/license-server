package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/models"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/antonioo83/license-server/internal/utils"
	"github.com/jinzhu/copier"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"strings"
	"time"
)

func InitCallbackCronJob(config config.Config, licenseRep interfaces.LicenseRepository) error {
	cronHandler := cron.New(cron.WithSeconds())
	cronHandler.AddFunc("1 * * * * *", func() {
		licenses, err := licenseRep.FindAllExpired(3, 50, 0)
		if err != nil {
			log.Printf("I can't get licenses for send callbacks: %v", err)
			return
		}
		callbackRequests, err := getLicenseCallbackRequests(licenses)
		if err != nil {
			log.Printf("I can't convert license models to callback requests: %v", err)
			return
		}
		runSendCallbackWorker(callbackRequests, licenseRep)
	})

	cronHandler.Start()

	return nil
}

type LicenseCallbackRequest struct {
	CustomerCode     string    `json:"customerId"`
	LicenseCode      string    `json:"licenseId"`
	ProductType      string    `json:"productType"`
	CallbackURL      string    `json:"callbackUrl"`
	Count            int       `json:"count"`
	LicenseKey       string    `json:"licenseKey"`
	ActivationAt     time.Time `json:"activationAt"`
	ExpirationAt     time.Time `json:"expirationAt"`
	Description      string    `json:"description"`
	LicenseId        int       `json:"-"`
	CallbackAttempts uint      `json:"-"`
	UserAuthToken    string    `json:"-"`
}

func getLicenseCallbackRequests(licences []models.Licence) ([]LicenseCallbackRequest, error) {
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

func runSendCallbackWorker(sendCallbacks []LicenseCallbackRequest, licenseRep interfaces.LicenseRepository) {
	for _, sendCallback := range sendCallbacks {
		go func() {
			httpStatus, err := sendRequest(sendCallback)
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
		}()
	}

	return
}

func sendRequest(callback LicenseCallbackRequest) (int, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	jsonRequest, err := getJSONRequest(callback)
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

func getJSONRequest(request LicenseCallbackRequest) ([]byte, error) {
	jsonResp, err := json.Marshal(request)
	if err != nil {
		return []byte(""), fmt.Errorf("i can't decode json request: %w", err)
	}

	return jsonResp, nil
}
