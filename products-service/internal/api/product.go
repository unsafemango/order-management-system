package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"yray.com/products-service/internal/models"
)

func (a *api) createProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeJSON(w, 400, map[string]string{"error": "bad json format"})
		return
	}

	if err := a.store.CreateProduct(r.Context(), &product); err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 201, product)
}

func (a *api) getProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	product, err := a.store.GetProduct(r.Context(), id)

	if err != nil {
		writeJSON(w, 404, map[string]string{"error": "product not found"})
		return
	}

	writeJSON(w, 200, product)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
