package events

import (
	"time"
)

type OrderCreated struct {
	EventID     string      `json:"eventId"`
	EventType   string      `json:"eventType"`
	Timestamp   time.Time   `json:"timestamp"`
	OrderID     string      `json:"orderId"`
	CustomerID  string      `json:"customerId"`
	TotalAmount float64     `json:"totalAmount"`
	Items       []OrderItem `json:"items"`
}

type OrderConfirmed struct {
	EventID   string    `json:"eventId"`
	EventType string    `json:"eventType"`
	Timestamp time.Time `json:"timestamp"`
	OrderID   string    `json:"orderId"`
}

type OrderCancelled struct {
	EventID   string    `json:"eventId"`
	EventType string    `json:"eventType"`
	Timestamp time.Time `json:"timestamp"`
	OrderID   string    `json:"orderId"`
	Reason    string    `json:"reason"`
}

type OrderItem struct {
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
