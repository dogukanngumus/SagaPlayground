package events

import (
	"time"
)

type PaymentRequested struct {
	EventID    string    `json:"eventId"`
	EventType  string    `json:"eventType"`
	Timestamp  time.Time `json:"timestamp"`
	OrderID    string    `json:"orderId"`
	CustomerID string    `json:"customerId"`
	Amount     float64   `json:"amount"`
	PaymentID  string    `json:"paymentId"`
}

type PaymentProcessed struct {
	EventID   string    `json:"eventId"`
	EventType string    `json:"eventType"`
	Timestamp time.Time `json:"timestamp"`
	OrderID   string    `json:"orderId"`
	PaymentID string    `json:"paymentId"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
}

type PaymentFailed struct {
	EventID   string    `json:"eventId"`
	EventType string    `json:"eventType"`
	Timestamp time.Time `json:"timestamp"`
	OrderID   string    `json:"orderId"`
	PaymentID string    `json:"paymentId"`
	Reason    string    `json:"reason"`
}
