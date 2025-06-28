package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"stock-service/handlers"
	"stock-service/services"

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

	stockService := services.NewStockService(db, broker)

	stockHandler := handlers.NewStockHandler(stockService)

	err = broker.ConsumeMessages("stock_events", func(body []byte) error {
		var paymentEvent map[string]interface{}
		if err := json.Unmarshal(body, &paymentEvent); err != nil {
			return err
		}

		if eventType, ok := paymentEvent["eventType"].(string); ok && eventType == "PaymentProcessed" {
			orderID := paymentEvent["orderId"].(string)

			log.Printf("Processing stock reservation for order: %s", orderID)

			items := []map[string]interface{}{
				{
					"productId": "product-1",
					"quantity":  2.0,
				},
			}

			if err := stockService.ReserveStock(orderID, items); err != nil {
				log.Printf("Error reserving stock: %v", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal("Failed to set up message consumer:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/products", stockHandler.GetAllProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", stockHandler.GetProduct).Methods("GET")

	port := getEnv("SERVICE_PORT", "5002")

	log.Printf("Stock service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
