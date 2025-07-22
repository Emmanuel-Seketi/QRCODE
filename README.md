# QR Code Management Platform - Backend

A comprehensive Go Fiber backend for managing QR codes with support for multiple QR code types, analytics, file uploads, and more.

## Features

- **Multiple QR Code Types**: Website, Search, Dynamic, Virtual Card, PDF, Social Media, Instagram, Images, App, Business, Event, 2D Barcode, Feedback, Rating, Email, Text, WiFi, SMS
- **CRUD Operations**: Create, Read, Update, Delete QR codes
- **File Upload Support**: PDFs, images, logos, photos
- **Analytics Tracking**: Scan insights with location, time, device tracking
- **Short URL Generation**: Dynamic QR codes with redirect handling
- **Tagging & Grouping**: Organize QR codes with tags and groups
- **Expiration & Deactivation**: Set expiration dates and deactivate QR codes
- **Design Customization**: Colors, logos, shapes
- **REST API**: Full REST API access

## Project Structure

```
QRFrontEnd/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── database/
│   │   └── database.go            # Database connection
│   ├── handler/
│   │   └── qr_handler.go          # QR code HTTP handlers
│   ├── model/
│   │   └── qr_code.go             # Data models and schemas
│   └── router/
│       └── router.go              # Route definitions
├── pkg/
│   ├── qrcode/
│   │   └── generator.go           # QR code generation utilities
│   └── shorturl/
│       └── generator.go           # Short URL generation utilities
├── .env                           # Environment variables
├── .gitignore                     # Git ignore rules
├── go.mod                         # Go module file
├── go.sum                         # Go dependencies checksum
└── README.md                      # This file
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (optional, currently using in-memory storage)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd QRFrontEnd
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:3000`

## API Endpoints

### QR Codes

- `POST /api/qr` - Create a new QR code
- `GET /api/qr/:id` - Get a QR code by ID
- `PUT /api/qr/:id` - Update a QR code
- `DELETE /api/qr/:id` - Delete a QR code

### Example Request

```bash
# Create a website QR code
curl -X POST http://localhost:3000/api/qr \
  -H "Content-Type: application/json" \
  -d '{
    "type": "website",
    "title": "My Website",
    "content": {
      "url": "https://example.com"
    },
    "analytics": true
  }'
```

## Configuration

The application uses environment variables for configuration. Key settings:

- `SERVER_PORT` - Server port (default: 3000)
- `SERVER_HOST` - Server host (default: localhost)
- `ENVIRONMENT` - Environment (development/production)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_NAME` - Database name
- `UPLOAD_PATH` - File upload directory
- `QR_CODE_SIZE` - Default QR code size
- `ANALYTICS_ENABLED` - Enable analytics tracking

## QR Code Types

The platform supports 19 different QR code types:

1. **Website** - Direct URL links
2. **Search** - Search engine queries
3. **Dynamic** - Updatable redirect URLs
4. **Virtual Card** - Contact information
5. **PDF** - PDF file links
6. **Social Media** - Social media profiles
7. **Instagram** - Instagram profiles
8. **Images** - Image galleries
9. **App** - App store links
10. **Business** - Business information
11. **Event** - Event details
12. **2D Barcode** - Custom data
13. **Feedback** - Feedback forms
14. **Rating** - Rating forms
15. **Email** - Email composition
16. **Text** - Plain text
17. **WiFi** - WiFi network details
18. **SMS** - SMS messages

## Development

### Adding New QR Code Types

1. Add the new type constant in `internal/model/qr_code.go`
2. Create the content structure for the new type
3. Update handlers to handle the new type
4. Add validation logic

### Database Migration

Currently using in-memory storage. To add database support:

1. Update `internal/database/database.go`
2. Add migration files
3. Update handlers to use database instead of in-memory storage

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License. 
