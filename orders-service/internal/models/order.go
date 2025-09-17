package models

import "time"

type Order struct {
	ID         int64       `json:"id"`
	CustomerID string      `json:"customer_id"`
	Status     string      `json:"status"` // NEW, PAID, CANCELLED, FULFILLED
	Items      []OrderItem `json:"items"`
	Total      int64       `json:"total"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
