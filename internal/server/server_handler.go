package server

import (
	"compress/flate"
	"github.com/antonioo83/license-server/config"
	handlers "github.com/antonioo83/license-server/internal/handlers"
	"github.com/antonioo83/license-server/internal/repositories/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"time"
)

type RouteParameters struct {
	Config                   config.Config
	UserRepository           interfaces.UserRepository
	UserActionRepository     interfaces.UserActionRepository
	UserPermissionRepository interfaces.UserPermissionRepository
}

func GetRouters(p RouteParameters) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)

	var params = handlers.UserRouteParameters{
		Config:               p.Config,
		UserRepository:       p.UserRepository,
		ActionRepository:     p.UserActionRepository,
		PermissionRepository: p.UserPermissionRepository,
	}
	r = getCreateUserRoute(r, params)
	r = getUpdateUserRoute(r, params)
	r = getDeleteUserRoute(r, params)
	r = getUserRoute(r, params)
	r = getUsersRoute(r, params)

	return r
}
