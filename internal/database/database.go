package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"qr_backend/ent"
	"qr_backend/internal/config"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// DB represents the Ent database client
var DB *ent.Client

// Connect initializes the database connection using Ent
func Connect(cfg *config.Config) error {
	var err error
	var dsn string
	var driver string

	switch cfg.Database.Type {
	case "sqlite":
		// Ensure data directory exists
		if err := os.MkdirAll(filepath.Dir(cfg.Database.Path), 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}

		dsn = fmt.Sprintf("file:%s?cache=shared&_fk=1", cfg.Database.Path)
		driver = "sqlite3"
		log.Printf("Using SQLite database at: %s", cfg.Database.Path)

	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)
		driver = "postgres"
		log.Printf("Using PostgreSQL database: %s@%s:%s/%s",
			cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	default:
		return fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	// Connect to database
	DB, err = ent.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to %s database: %w", cfg.Database.Type, err)
	}

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create/update database schema
	if err = DB.Schema.Create(ctx); err != nil {
		DB.Close()
		return fmt.Errorf("failed to create %s schema: %w", cfg.Database.Type, err)
	}

	log.Printf("%s database connected and schema created successfully", cfg.Database.Type)
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
