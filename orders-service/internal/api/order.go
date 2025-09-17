package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"yray.com/orders-service/internal/models"
)

func (a *api) createOrder(w http.ResponseWriter, r *http.Request) {
	var req createOrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "bad json format"})
	}
	if len(req.Items) == 0 {
		writeJSON(w, 400, map[string]string{"error": "no items selected"})
		return
	}

	// pricing fetch + reserve stock in products-service
	var total int64
	for i, it := range req.Items {
		price, err := a.fetchPrice(it.ProductID)
		if err != nil {
			writeJSON(w, 400, map[string]string{
				"error": err.Error(),
			})
			return
		}
		req.Items[i].Price = price
		total += price * it.Qty
		ok, err := a.reserve(it.ProductID, it.Qty)
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		if !ok {
			writeJSON(w, 409, map[string]string{"error": "insufficient stock"})
			return
		}
	}

	o := &models.Order{
		CustomerID: req.CustomerID,
		Status:     "NEW",
		Items:      req.Items,
		Total:      total,
	}
	id, err := a.store.CreateOrder(r.Context(), o)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, map[string]any{"id": id, "total": total, "status": "NEW"})
}

func (a *api) getOrder(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	o, err := a.store.GetOrder(r.Context(), id)
	if err != nil {
		writeJSON(w, 404, map[string]string{"error": "order not found"})
		return
	}

	writeJSON(w, 200, o)
}

func (a *api) updateStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var req updateStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "bad json format"})
		return
	}
	if err := a.store.UpdateStatus(r.Context(), id, req.Status); err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]string{"status": req.Status})
}
