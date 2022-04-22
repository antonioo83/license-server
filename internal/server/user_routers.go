package server

import (
	"github.com/antonioo83/license-server/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func getCreateUserRoute(r *chi.Mux, params handlers.UserRouteParameters) *chi.Mux {
	r.Post("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreatedUserResponse(r, w, params)
	})

	return r
}

func getUpdateUserRoute(r *chi.Mux, params handlers.UserRouteParameters) *chi.Mux {
	r.Put("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreatedUserResponse(r, w, params)
	})

	return r
}

func getDeleteUserRoute(r *chi.Mux, params handlers.UserRouteParameters) *chi.Mux {
	r.Delete("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDeletedUserResponse(r, w, params)
	})

	return r
}

func getUserRoute(r *chi.Mux, params handlers.UserRouteParameters) *chi.Mux {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreatedUserResponse(r, w, params)
	})

	return r
}

func getUsersRoute(r *chi.Mux, params handlers.UserRouteParameters) *chi.Mux {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreatedUserResponse(r, w, params)
	})

	return r
}
