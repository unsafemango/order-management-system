package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"yray.com/products-service/internal/repo"
)

type api struct {
	store *repo.Store
}

func NewRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	api := &api{
		store: repo.NewStore(db),
	}

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", api.healthCheckHandler)
		r.Post("/products", api.createProduct)
		r.Get("/products/{id}", api.getProduct)
		r.Post("/products/{id}/reserve", api.reserve)
		r.Post("/products/{id}/release", api.release)
	})

	return r
}
