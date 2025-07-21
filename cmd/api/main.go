package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"qr_backend/internal/config"
	"qr_backend/internal/database"
	"qr_backend/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Add this import
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database connection
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Setup graceful shutdown
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize the template engine with absolute path for robustness
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current working directory:", err)
	}
	engine := html.New(filepath.Join(cwd, "templates"), ".html")
	if err := engine.Load(); err != nil {
		log.Fatal("Failed to load templates:", err)
	}

	// Create Fiber app with configuration and template engine
	app := fiber.New(fiber.Config{
		AppName: "QR Code Management Platform",
		Views:   engine,
	})

	// Add Fiber logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip} - ${method} ${path} - ${status}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Add CORS middleware with proper configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://127.0.0.1:5173,http://192.168.198.78:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
	}))

	// Setup routes
	router.SetupRoutes(app)

	// Rest of your existing code...
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s in %s mode", serverAddr, cfg.Server.Environment)

	if err := app.Listen(serverAddr); err != nil {
		log.Printf("Server is shutting down: %v", err)
	}
}
