package main

import (
	"fmt"
	"log"
	"plantbased-backend/config"
	"plantbased-backend/database"
	"plantbased-backend/utils"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations first
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(cfg.AdminPassword)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Insert admin
	_, err = db.Exec(`
		INSERT INTO admins (email, password_hash, full_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO NOTHING
	`, cfg.AdminEmail, hashedPassword, "System Administrator", true)

	if err != nil {
		log.Fatal("Failed to create admin:", err)
	}

	fmt.Println("✅ Admin user created successfully!")
	fmt.Printf("📧 Email: %s\n", cfg.AdminEmail)
	fmt.Println("🔐 Password: (check your .env file)")
}