package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Server
	Port string
	Env  string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// Admin
	AdminEmail    string
	AdminPassword string

	// Cloudinary
	CloudinaryCloudName   string
	CloudinaryAPIKey      string
	CloudinaryAPISecret   string
	CloudinaryUploadFolder string
}

var AppConfig *Config

// LoadConfig loads environment variables into Config struct
func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	AppConfig = &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "plantbased_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Server
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		// JWT
		JWTSecret:      getEnv("JWT_SECRET", "change-this-secret"),
		JWTExpiryHours: jwtExpiry,

		// Admin
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@plantbased.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
	}

	return AppConfig
}

// getEnv gets environment variable with fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}