package models

import (
	"time"
)

type Payment struct {
	ID         string    `json:"id" db:"id"`
	OrderID    string    `json:"orderId" db:"order_id"`
	CustomerID string    `json:"customerId" db:"customer_id"`
	Amount     float64   `json:"amount" db:"amount"`
	Status     string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type OutboxMessage struct {
	ID          string     `json:"id" db:"id"`
	EventType   string     `json:"eventType" db:"event_type"`
	EventData   string     `json:"eventData" db:"event_data"`
	Exchange    string     `json:"exchange" db:"exchange"`
	RoutingKey  string     `json:"routingKey" db:"routing_key"`
	IsProcessed bool       `json:"isProcessed" db:"is_processed"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	ProcessedAt *time.Time `json:"processedAt" db:"processed_at"`
}
