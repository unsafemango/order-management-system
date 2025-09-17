package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"yray.com/orders-service/internal/repo"
)

type api struct {
	store       *repo.Store
	productsURL string
}

func NewRouter(db *pgxpool.Pool, productsURL string) http.Handler {
	r := chi.NewRouter()
	api := &api{
		store:       repo.NewStore(db),
		productsURL: productsURL,
	}

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", api.healthCheckHandler)
		r.Post("/orders", api.createOrder)
		r.Get("/orders/{id}", api.getOrder)
		r.Patch("/orders/{id}/status", api.updateStatus)
	})

	return r
}
