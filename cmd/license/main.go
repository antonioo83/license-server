package main

import (
	"context"
	"fmt"
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/handlers"
	authFactory "github.com/antonioo83/license-server/internal/handlers/auth/factory"
	"github.com/antonioo83/license-server/internal/repositories/factory"
	"github.com/antonioo83/license-server/internal/server"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func main() {
	config := config.GetConfigSettings()
	var pool *pgxpool.Pool
	context := context.Background()
	pool, _ = pgxpool.Connect(context, config.DatabaseDsn)
	defer pool.Close()

	licenseRepository := factory.NewLicenseRepository(context, pool)
	err := handlers.InitCallbackCronJob(config.Callback, licenseRepository)
	if err != nil {
		fmt.Println("i can't run send callback job: " + err.Error())
	}

	userPermissionRepository := factory.NewUserPermissionRepository(context, pool)
	userRepository := factory.NewUserRepository(context, pool, userPermissionRepository)
	userAuthHandler := authFactory.NewUserAuthHandler(userRepository, config)
	routeParameters :=
		server.RouteParameters{
			Config:                   config,
			UserRepository:           userRepository,
			UserActionRepository:     factory.NewUserActionRepository(context, pool),
			UserPermissionRepository: userPermissionRepository,
		}
	licenseRouteParameters :=
		handlers.LicenseRouteParameters{
			Config:             config,
			CustomerRepository: factory.NewCustomerRepository(context, pool, licenseRepository),
			LicenseRepository:  licenseRepository,
		}

	handler := server.GetRouters(userAuthHandler, routeParameters, licenseRouteParameters)
	log.Fatal(http.ListenAndServe(config.ServerAddress, handler))
}
