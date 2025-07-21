package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config struct holds the application configuration
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	Upload    UploadConfig
	QRCode    QRCodeConfig
	Analytics AnalyticsConfig
	Redis     RedisConfig
	External  ExternalConfig
	Logging   LoggingConfig
}

type ServerConfig struct {
	Port        string
	Host        string
	Environment string
}

type DatabaseConfig struct {
	Type     string // "sqlite" or "postgres"
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	Path     string // For SQLite
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type UploadConfig struct {
	MaxSize      int64
	Path         string
	AllowedTypes []string
}

type QRCodeConfig struct {
	Size   int
	Level  string
	Margin int
}

type AnalyticsConfig struct {
	Enabled       bool
	RetentionDays int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type ExternalConfig struct {
	ShortURLDomain string
	CORSOrigin     string
}

type LoggingConfig struct {
	Level string
	File  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "3000"),
			Host:        getEnv("SERVER_HOST", "localhost"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "sqlite"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "qr_platform"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			Path:     getEnv("DB_PATH", "./data/qr_platform.db"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your_jwt_secret_here"),
			Expiry: getEnvDuration("JWT_EXPIRY", 24*time.Hour),
		},
		Upload: UploadConfig{
			MaxSize:      getEnvInt64("UPLOAD_MAX_SIZE", 10485760), // 10MB
			Path:         getEnv("UPLOAD_PATH", "./uploads"),
			AllowedTypes: getEnvSlice("ALLOWED_FILE_TYPES", []string{"pdf", "jpg", "jpeg", "png", "gif", "svg"}),
		},
		QRCode: QRCodeConfig{
			Size:   getEnvInt("QR_CODE_SIZE", 256),
			Level:  getEnv("QR_CODE_LEVEL", "M"),
			Margin: getEnvInt("QR_CODE_MARGIN", 1),
		},
		Analytics: AnalyticsConfig{
			Enabled:       getEnvBool("ANALYTICS_ENABLED", true),
			RetentionDays: getEnvInt("ANALYTICS_RETENTION_DAYS", 365),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		External: ExternalConfig{
			ShortURLDomain: getEnv("SHORT_URL_DOMAIN", "qr.yourdomain.com"),
			CORSOrigin:     getEnv("CORS_ORIGIN", "http://localhost:3000"),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "./logs/app.log"),
		},
	}

	return config, nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated values
		return []string{value} // You might want to split by comma here
	}
	return defaultValue
}
