package server

import (
	"github.com/antonioo83/license-server/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func getReplaceLicenseRoute(r *chi.Mux, params handlers.LicenseRouteParameters) *chi.Mux {
	r.Post("/api/v1/licenses/replace", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetReplacedLicenseResponse(r, w, params)
	})

	return r
}

func getDeleteLicenseRoute(r *chi.Mux, params handlers.LicenseRouteParameters) *chi.Mux {
	r.Delete("/api/v1/licenses", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDeletedLicenseResponse(r, w, params)
	})

	return r
}

func getLicenseRoute(r *chi.Mux, params handlers.LicenseRouteParameters) *chi.Mux {
	r.Get("/api/v1/licenses", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetLicenseResponse(r, w, params)
	})

	return r
}

func getLicensesRoute(r *chi.Mux, params handlers.LicenseRouteParameters) *chi.Mux {
	r.Get("/api/v1/licenses", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetLicensesResponse(r, w, params)
	})

	return r
}
