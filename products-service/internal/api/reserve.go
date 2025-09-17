package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (a *api) reserve(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var req struct {
		Qty int64 `json:"qty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "bad json format"})
		return
	}

	reserve, err := a.store.Reserve(context.Background(), id, req.Qty)

	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	if !reserve.OK {
		writeJSON(w, 409, reserve)
		return
	}

	writeJSON(w, 200, reserve)
}

func (a *api) release(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var req struct {
		Qty int64 `json:"qty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "bad json format"})
		return
	}

	if err := a.store.Release(r.Context(), id, req.Qty); err != nil {
		writeJSON(w, 404, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]string{"status": "released"})
}
