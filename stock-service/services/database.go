package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DatabaseService struct {
	DB *sql.DB
}

func NewDatabaseService() (*DatabaseService, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASS", "postgres123")
	dbname := getEnv("DB_NAME", "saga_demo")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, err
	}

	// Initialize sample products
	if err := initializeProducts(db); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return &DatabaseService{DB: db}, nil
}

func createTables(db *sql.DB) error {
	// Create products table
	productsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		quantity INTEGER NOT NULL,
		price DECIMAL(10,2) NOT NULL
	);`

	// Create stock_reservations table
	reservationsTable := `
	CREATE TABLE IF NOT EXISTS stock_reservations (
		id VARCHAR(255) PRIMARY KEY,
		order_id VARCHAR(255) NOT NULL,
		product_id VARCHAR(255) NOT NULL,
		quantity INTEGER NOT NULL,
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Create outbox_messages table
	outboxTable := `
	CREATE TABLE IF NOT EXISTS outbox_messages (
		id VARCHAR(255) PRIMARY KEY,
		event_type VARCHAR(255) NOT NULL,
		event_data TEXT NOT NULL,
		exchange VARCHAR(255) NOT NULL,
		routing_key VARCHAR(255) NOT NULL,
		is_processed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		processed_at TIMESTAMP NULL
	);`

	if _, err := db.Exec(productsTable); err != nil {
		return err
	}

	if _, err := db.Exec(reservationsTable); err != nil {
		return err
	}

	if _, err := db.Exec(outboxTable); err != nil {
		return err
	}

	return nil
}

func initializeProducts(db *sql.DB) error {
	// Check if products already exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Products already exist
	}

	// Insert sample products
	products := []struct {
		id       string
		name     string
		quantity int
		price    float64
	}{
		{"product-1", "Laptop", 10, 999.99},
		{"product-2", "Mouse", 50, 29.99},
		{"product-3", "Keyboard", 30, 79.99},
	}

	for _, p := range products {
		_, err := db.Exec("INSERT INTO products (id, name, quantity, price) VALUES ($1, $2, $3, $4)",
			p.id, p.name, p.quantity, p.price)
		if err != nil {
			return err
		}
	}

	log.Println("Sample products initialized")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
