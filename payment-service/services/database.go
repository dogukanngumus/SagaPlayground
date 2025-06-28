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

	log.Println("Database connected successfully")
	return &DatabaseService{DB: db}, nil
}

func createTables(db *sql.DB) error {
	// Create payments table
	paymentsTable := `
	CREATE TABLE IF NOT EXISTS payments (
		id VARCHAR(255) PRIMARY KEY,
		order_id VARCHAR(255) NOT NULL,
		customer_id VARCHAR(255) NOT NULL,
		amount DECIMAL(10,2) NOT NULL,
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

	if _, err := db.Exec(paymentsTable); err != nil {
		return err
	}

	if _, err := db.Exec(outboxTable); err != nil {
		return err
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
