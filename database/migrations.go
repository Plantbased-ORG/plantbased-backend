package database

import (
	"database/sql"
	"fmt"
)

// RunMigrations creates necessary tables
func RunMigrations(db *sql.DB) error {
	// Create admins table
	createAdminsTable := `
	CREATE TABLE IF NOT EXISTS admins (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		full_name VARCHAR(255),
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createAdminsTable); err != nil {
		return fmt.Errorf("failed to create admins table: %w", err)
	}

	// Create index on email
	createEmailIndex := `
	CREATE INDEX IF NOT EXISTS idx_admins_email ON admins(email);
	`

	if _, err := db.Exec(createEmailIndex); err != nil {
		return fmt.Errorf("failed to create email index: %w", err)
	}

	// Create programs table
	createProgramsTable := `
	CREATE TABLE IF NOT EXISTS programs (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		short_description TEXT NOT NULL,
		main_image_public_id VARCHAR(255) NOT NULL,
		main_image_url TEXT NOT NULL,
		intro_description TEXT NOT NULL,
		main_content_image_public_id VARCHAR(255) NOT NULL,
		main_content_image_url TEXT NOT NULL,
		what_causes TEXT NOT NULL,
		what_causes_image_public_id VARCHAR(255) NOT NULL,
		what_causes_image_url TEXT NOT NULL,
		health_risks TEXT NOT NULL,
		health_risks_image_public_id VARCHAR(255) NOT NULL,
		health_risks_image_url TEXT NOT NULL,
		strategies TEXT NOT NULL,
		strategies_image_public_id VARCHAR(255) NOT NULL,
		strategies_image_url TEXT NOT NULL,
		conclusion TEXT NOT NULL,
		conclusion_image_public_id VARCHAR(255) NOT NULL,
		conclusion_image_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createProgramsTable); err != nil {
		return fmt.Errorf("failed to create programs table: %w", err)
	}

	// Create program_pricing_plans table
	createPricingPlansTable := `
	CREATE TABLE IF NOT EXISTS program_pricing_plans (
		id SERIAL PRIMARY KEY,
		program_id INTEGER NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		subtitle TEXT NOT NULL,
		price VARCHAR(50) NOT NULL,
		features JSONB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createPricingPlansTable); err != nil {
		return fmt.Errorf("failed to create program_pricing_plans table: %w", err)
	}

	// Create index on program_id
	createProgramIDIndex := `
	CREATE INDEX IF NOT EXISTS idx_program_pricing_plans_program_id ON program_pricing_plans(program_id);
	`

	if _, err := db.Exec(createProgramIDIndex); err != nil {
		return fmt.Errorf("failed to create program_id index: %w", err)
	}

	// Create testimonials table
	createTestimonialsTable := `
	CREATE TABLE IF NOT EXISTS testimonials (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		location VARCHAR(255) NOT NULL,
		review TEXT NOT NULL,
		avatar TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createTestimonialsTable); err != nil {
		return fmt.Errorf("failed to create testimonials table: %w", err)
	}

	return nil
}