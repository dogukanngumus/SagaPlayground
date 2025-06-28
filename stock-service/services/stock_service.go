package services

import (
	"log"
	"time"

	"stock-service/models"

	"github.com/google/uuid"
)

type StockService struct {
	db     *DatabaseService
	broker *MessageBrokerService
}

func NewStockService(db *DatabaseService, broker *MessageBrokerService) *StockService {
	return &StockService{
		db:     db,
		broker: broker,
	}
}

func (s *StockService) ReserveStock(orderID string, items []map[string]interface{}) error {
	for _, item := range items {
		productID := item["productId"].(string)
		quantity := int(item["quantity"].(float64))

		var currentStock int
		err := s.db.DB.QueryRow("SELECT quantity FROM products WHERE id = $1", productID).Scan(&currentStock)
		if err != nil {
			log.Printf("Product not found: %s", productID)
			return err
		}

		if currentStock < quantity {
			log.Printf("Insufficient stock for product %s: requested %d, available %d", productID, quantity, currentStock)
			return s.publishStockReservationFailed(orderID, "Insufficient stock")
		}

		reservationID := uuid.New().String()
		reservation := &models.StockReservation{
			ID:        reservationID,
			OrderID:   orderID,
			ProductID: productID,
			Quantity:  quantity,
			Status:    "Reserved",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.saveReservation(reservation); err != nil {
			return err
		}

		if err := s.updateProductStock(productID, currentStock-quantity); err != nil {
			return err
		}
	}

	event := map[string]interface{}{
		"eventId":   uuid.New().String(),
		"eventType": "StockReserved",
		"timestamp": time.Now(),
		"orderId":   orderID,
		"items":     items,
	}

	if err := s.broker.PublishMessage("saga.events", "stock.reserved", event); err != nil {
		log.Printf("Error publishing StockReserved event: %v", err)
		return err
	}

	log.Printf("Stock reserved successfully for order: %s", orderID)
	return nil
}

func (s *StockService) saveReservation(reservation *models.StockReservation) error {
	query := `
		INSERT INTO stock_reservations (id, order_id, product_id, quantity, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.db.DB.Exec(query, reservation.ID, reservation.OrderID, reservation.ProductID,
		reservation.Quantity, reservation.Status, reservation.CreatedAt, reservation.UpdatedAt)
	return err
}

func (s *StockService) updateProductStock(productID string, newQuantity int) error {
	query := `UPDATE products SET quantity = $1 WHERE id = $2`
	_, err := s.db.DB.Exec(query, newQuantity, productID)
	return err
}

func (s *StockService) publishStockReservationFailed(orderID, reason string) error {
	event := map[string]interface{}{
		"eventId":   uuid.New().String(),
		"eventType": "StockReservationFailed",
		"timestamp": time.Now(),
		"orderId":   orderID,
		"reason":    reason,
	}

	return s.broker.PublishMessage("saga.events", "stock.reservation.failed", event)
}

func (s *StockService) GetProduct(productID string) (*models.Product, error) {
	query := `SELECT id, name, quantity, price FROM products WHERE id = $1`

	var product models.Product
	err := s.db.DB.QueryRow(query, productID).Scan(&product.ID, &product.Name, &product.Quantity, &product.Price)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *StockService) GetAllProducts() ([]models.Product, error) {
	query := `SELECT id, name, quantity, price FROM products ORDER BY name`

	rows, err := s.db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
