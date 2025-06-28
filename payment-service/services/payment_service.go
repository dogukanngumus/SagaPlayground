package services

import (
	"log"
	"time"

	"payment-service/models"

	"github.com/google/uuid"
)

type PaymentService struct {
	db     *DatabaseService
	broker *MessageBrokerService
}

func NewPaymentService(db *DatabaseService, broker *MessageBrokerService) *PaymentService {
	return &PaymentService{
		db:     db,
		broker: broker,
	}
}

func (p *PaymentService) ProcessPaymentRequest(orderID, customerID string, amount float64) error {
	paymentID := uuid.New().String()
	payment := &models.Payment{
		ID:         paymentID,
		OrderID:    orderID,
		CustomerID: customerID,
		Amount:     amount,
		Status:     "Processing",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := p.savePayment(payment); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	payment.Status = "Completed"
	payment.UpdatedAt = time.Now()

	if err := p.updatePayment(payment); err != nil {
		return err
	}

	event := map[string]interface{}{
		"eventId":   uuid.New().String(),
		"eventType": "PaymentProcessed",
		"timestamp": time.Now(),
		"orderId":   orderID,
		"paymentId": paymentID,
		"amount":    amount,
		"status":    "Completed",
	}

	if err := p.broker.PublishMessage("saga.events", "payment.processed", event); err != nil {
		log.Printf("Error publishing PaymentProcessed event: %v", err)
		return err
	}

	log.Printf("Payment processed successfully: %s", paymentID)
	return nil
}

func (p *PaymentService) savePayment(payment *models.Payment) error {
	query := `
		INSERT INTO payments (id, order_id, customer_id, amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := p.db.DB.Exec(query, payment.ID, payment.OrderID, payment.CustomerID,
		payment.Amount, payment.Status, payment.CreatedAt, payment.UpdatedAt)
	return err
}

func (p *PaymentService) updatePayment(payment *models.Payment) error {
	query := `
		UPDATE payments 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := p.db.DB.Exec(query, payment.Status, payment.UpdatedAt, payment.ID)
	return err
}

func (p *PaymentService) GetPayment(paymentID string) (*models.Payment, error) {
	query := `SELECT id, order_id, customer_id, amount, status, created_at, updated_at 
			  FROM payments WHERE id = $1`

	var payment models.Payment
	err := p.db.DB.QueryRow(query, paymentID).Scan(
		&payment.ID, &payment.OrderID, &payment.CustomerID, &payment.Amount,
		&payment.Status, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (p *PaymentService) GetAllPayments() ([]models.Payment, error) {
	query := `SELECT id, order_id, customer_id, amount, status, created_at, updated_at 
			  FROM payments ORDER BY created_at DESC`

	rows, err := p.db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(&payment.ID, &payment.OrderID, &payment.CustomerID, &payment.Amount,
			&payment.Status, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}
