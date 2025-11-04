package database

import (
	"database/sql"
	"fmt"
	"plantbased-backend/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes and returns a database connection
func InitDB(cfg *config.Config) (*sql.DB, error) {
	// Build connection string with SSL mode
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	DB = db
	fmt.Println("✅ Connected to database successfully")
	return db, nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
}

// testing something here