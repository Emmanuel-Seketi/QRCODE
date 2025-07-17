package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"qr_backend/ent"
	"qr_backend/internal/config"

	_ "github.com/lib/pq"
)

// DB represents the Ent database client
var DB *ent.Client

// Connect initializes the database connection using Ent
func Connect(cfg *config.Config) error {
	var err error
	
	// PostgreSQL connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
	
	// Connect to PostgreSQL
	DB, err = ent.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	
	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Create/update database schema
	if err = DB.Schema.Create(ctx); err != nil {
		DB.Close()
		return fmt.Errorf("failed to create PostgreSQL schema: %w", err)
	}
	
	log.Println("PostgreSQL database connected and schema created successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
