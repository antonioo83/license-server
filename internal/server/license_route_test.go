package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/handlers"
	"github.com/antonioo83/license-server/internal/repositories/factory"
	"github.com/antonioo83/license-server/internal/utils"
	"github.com/bxcodec/faker/v3"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type LicenseTest struct {
	LicenseId    string `faker:"uuid_hyphenated" json:"licenseId,omitempty"`
	ProductType  string `faker:"oneof: courier, waiter, pechka54, solo" json:"productType,omitempty"`
	CallbackURL  string `faker:"url" json:"callbackUrl,omitempty"`
	Count        int    `faker:"boundary_start=1, boundary_end=100" json:"count,omitempty"`
	LicenseKey   string `faker:"uuid_hyphenated" json:"licenseKey,omitempty"`
	ActivationAt string `faker:"timestamp" json:"activationAt,omitempty"`
	ExpirationAt string `faker:"timestamp" json:"expirationAt,omitempty"`
	Description  string `faker:"len=256" json:"description,omitempty"`
}

type CustomerTest struct {
	CustomerId  string         `faker:"uuid_hyphenated" json:"customerId,omitempty"`
	Type        string         `faker:"oneof: device, service" json:"type,omitempty"`
	Inn         string         `faker:"-" json:"inn,omitempty"`
	Title       string         `faker:"username" json:"title,omitempty"`
	Description string         `faker:"len=256" json:"description,omitempty"`
	Licenses    [1]LicenseTest `json:"Licenses,omitempty"`
}

func TestCRUDLicenseRouters(t *testing.T) {
	licenseTests := []struct {
		url     string
		request CustomerTest
	}{
		{
			url: "/api/v1/licenses/replace",
			request: CustomerTest{
				CustomerId:  "",
				Type:        "",
				Inn:         "",
				Title:       "",
				Description: "",
				Licenses: [1]LicenseTest{{
					LicenseId:    "",
					ProductType:  "",
					CallbackURL:  "",
					Count:        0,
					LicenseKey:   "",
					ActivationAt: "",
					ExpirationAt: "",
					Description:  "",
				}},
			},
		},
	}

	var pool *pgxpool.Pool
	context := context.Background()
	config := config.GetConfigSettings()

	pool, _ = pgxpool.Connect(context, config.DatabaseDsn)
	defer pool.Close()
	userPermissionRepository := factory.NewUserPermissionRepository(context, pool)
	routeParameters :=
		RouteParameters{
			Config:                   config,
			UserRepository:           factory.NewUserRepository(context, pool, userPermissionRepository),
			UserActionRepository:     factory.NewUserActionRepository(context, pool),
			UserPermissionRepository: userPermissionRepository,
		}

	licenseRepository := factory.NewLicenseRepository(context, pool)
	licenseRouteParameters :=
		handlers.LicenseRouteParameters{
			Config:             config,
			CustomerRepository: factory.NewCustomerRepository(context, pool, licenseRepository),
			LicenseRepository:  licenseRepository,
		}
	r := GetRouters(routeParameters, licenseRouteParameters)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range licenseTests {
		err := faker.FakeData(&tt.request)
		if err != nil {
			log.Fatal(err)
		}
		tt.request.Inn = "1234567890"

		request, err := getCustomerRequest(tt.request)
		assert.NoError(t, err)

		jsonRequest := getPostLicenseRequest(t, ts, tt.url, strings.NewReader(string(request)), config.Auth.AdminAuthToken)
		resp, _ := sendLicenseRequest(t, jsonRequest)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}

}

func getCustomerRequest(request CustomerTest) ([]byte, error) {
	jsonResp, err := json.Marshal(request)
	if err != nil {
		return []byte(""), fmt.Errorf("i can't decode json request: %w", err)
	}

	return jsonResp, nil
}

func getPostLicenseRequest(t *testing.T, ts *httptest.Server, path string, body io.Reader, token string) *http.Request {
	req, err := http.NewRequest("POST", ts.URL+path, body)
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	require.NoError(t, err)

	return req
}

func sendLicenseRequest(t *testing.T, req *http.Request) (*http.Response, string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer utils.ResourceClose(resp.Body)

	return resp, string(respBody)
}
