package api

import "yray.com/orders-service/internal/models"

type createOrderReq struct {
	CustomerID string             `json:"customer_id"`
	Items      []models.OrderItem `json:"items"`
}

type priceLookup struct {
	Price int64 `json:"price"`
}

type reserveReq struct {
	Qty int64 `json:"qty"`
}

type reserveResp struct {
	OK        bool   `json:"ok"`
	Available int64  `json:"available"`
	Message   string `json:"message"`
}

type updateStatusReq struct {
	Status string `json:"status"`
}
