package models

type OrderItem struct {
	ProductID int64 `json:"product_id"`
	Qty       int64 `json:"qty"`
	Price     int64 `json:"price"`
}
