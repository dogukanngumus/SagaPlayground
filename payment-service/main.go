package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"payment-service/handlers"
	"payment-service/services"

	"github.com/gorilla/mux"
)

func main() {
	db, err := services.NewDatabaseService()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.DB.Close()

	broker, err := services.NewMessageBrokerService()
	if err != nil {
		log.Fatal("Failed to connect to message broker:", err)
	}
	defer broker.Close()

	paymentService := services.NewPaymentService(db, broker)

	paymentHandler := handlers.NewPaymentHandler(paymentService)

	err = broker.ConsumeMessages("payment_events", func(body []byte) error {
		var orderEvent map[string]interface{}
		if err := json.Unmarshal(body, &orderEvent); err != nil {
			return err
		}

		if eventType, ok := orderEvent["eventType"].(string); ok && eventType == "OrderCreated" {
			orderID := orderEvent["orderId"].(string)
			customerID := orderEvent["customerId"].(string)
			totalAmount := orderEvent["totalAmount"].(float64)

			log.Printf("Processing payment for order: %s, amount: %.2f", orderID, totalAmount)

			if err := paymentService.ProcessPaymentRequest(orderID, customerID, totalAmount); err != nil {
				log.Printf("Error processing payment: %v", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal("Failed to set up message consumer:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/payments", paymentHandler.GetAllPayments).Methods("GET")
	r.HandleFunc("/api/payments/{id}", paymentHandler.GetPayment).Methods("GET")

	port := getEnv("SERVICE_PORT", "5001")

	log.Printf("Payment service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
