// ReserveResult is returned when reserving a stock for an order. If OK is false, not enough stock is available
package models

type ReserveResult struct {
	OK        bool   `json:"ok"`
	Available int64  `json:"available"`
	Message   string `json:"message"`
}
