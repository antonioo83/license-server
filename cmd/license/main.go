package main

import (
	"context"
	"github.com/antonioo83/license-server/config"
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
	userPermissionRepository := factory.NewUserPermissionRepository(context, pool)
	routeParameters :=
		server.RouteParameters{
			Config:                   config,
			UserRepository:           factory.NewUserRepository(context, pool, userPermissionRepository),
			UserActionRepository:     factory.NewUserActionRepository(context, pool),
			UserPermissionRepository: userPermissionRepository,
		}
	handler := server.GetRouters(routeParameters)
	log.Fatal(http.ListenAndServe(config.ServerAddress, handler))
}
