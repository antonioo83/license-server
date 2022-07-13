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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Product struct {
	Type        string    `faker:"oneof: courier, waiter, pechka54, solo" json:"type,omitempty"`
	Permissions [4]string `faker:"-" json:"permissions,omitempty"`
}

type RequestTest struct {
	UserId      string     `faker:"uuid_hyphenated" json:"userId,omitempty"`
	Role        string     `faker:"oneof: service,device" json:"role,omitempty"`
	Title       string     `faker:"username" json:"title,omitempty"`
	Description string     `faker:"len=256" json:"description,omitempty"`
	Products    [1]Product `json:"products,omitempty"`
}

func TestGetRouters(t *testing.T) {
	userTests := []struct {
		url     string
		request RequestTest
	}{
		{
			url: "/api/v1/users",
			request: RequestTest{
				UserId: "",
				Role:   "",
				Title:  "",
				Products: [1]Product{{
					Type:        "",
					Permissions: [4]string{"create", "update", "delete", "get"},
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

	for _, tt := range userTests {
		err := faker.FakeData(&tt.request)
		assert.NoError(t, err)
		tt.request.Products[0].Permissions = [4]string{"create", "update", "delete", "get"}

		request, err := getJSONRequest(tt.request)
		assert.NoError(t, err)

		jsonRequest := getPostRequest(t, ts, tt.url, strings.NewReader(string(request)), config.Auth.AdminAuthToken)
		resp, _ := sendRequest(t, jsonRequest)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
		resp.Body.Close()
		assert.NoError(t, err)
	}
}

func getJSONRequest(request RequestTest) ([]byte, error) {
	jsonResp, err := json.Marshal(request)
	if err != nil {
		return []byte(""), fmt.Errorf("i can't decode json request: %w", err)
	}

	return jsonResp, nil
}

func getPostRequest(t *testing.T, ts *httptest.Server, path string, body io.Reader, token string) *http.Request {
	req, err := http.NewRequest("POST", ts.URL+path, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	require.NoError(t, err)

	return req
}

func sendRequest(t *testing.T, req *http.Request) (*http.Response, string) {
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
