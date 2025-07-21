package router

import (
	"qr_backend/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// QR Code routes
	qr := api.Group("/qr")
	qr.Get("/", handler.ListQRCodes)                     // List QR codes with pagination
	qr.Post("/", handler.CreateQRCode)                   // Create a new QR code
	qr.Post("/pdf", handler.CreatePDFQRCode)             // Create PDF QR code with file upload
	qr.Get("/:id", handler.GetQRCode)                    // Get a QR code by ID
	qr.Put("/:id", handler.UpdateQRCode)                 // Update a QR code
	qr.Delete("/:id", handler.DeleteQRCode)              // Delete a QR code
	qr.Delete("/", handler.BulkDeleteQRCodes)            // Bulk delete all QR codes
	qr.Get("/:id/download", handler.DownloadQRCode)      // Download QR code image
	qr.Get("/:id/analytics", handler.GetQRCodeAnalytics) // Get QR code analytics

	// File upload routes
	api.Post("/upload", handler.UploadFile) // Upload files

	// Scan/redirect routes (outside API group for clean URLs)
	app.Get("/scan/:shortcode", handler.ScanQRCode) // QR code scanning and redirection
	app.Get("/qr/:id", handler.GetStaticQRContent)  // Static QR content display

	// Static file serving for uploads
	app.Static("/uploads", "./uploads")
}
