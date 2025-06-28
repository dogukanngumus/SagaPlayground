package events

import (
	"time"
)

type StockReserved struct {
	EventID   string      `json:"eventId"`
	EventType string      `json:"eventType"`
	Timestamp time.Time   `json:"timestamp"`
	OrderID   string      `json:"orderId"`
	Items     []StockItem `json:"items"`
}

type StockReservationFailed struct {
	EventID   string    `json:"eventId"`
	EventType string    `json:"eventType"`
	Timestamp time.Time `json:"timestamp"`
	OrderID   string    `json:"orderId"`
	Reason    string    `json:"reason"`
}

type StockItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
