package main

import (
	"log"
	"plantbased-backend/config"
	"plantbased-backend/database"
	"plantbased-backend/middleware"
	"plantbased-backend/routes"
	"plantbased-backend/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting PlantBased Backend...")

	// Load configuration
	log.Println("Loading configuration...")
	cfg := config.LoadConfig()
	log.Printf("✓ Configuration loaded (Port: %s, Env: %s)", cfg.Port, cfg.Env)

	// Connect to database
	log.Println("Connecting to database...")
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	log.Println("✓ Database connected successfully")

	// Run migrations
	log.Println("Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("✓ Migrations completed successfully")

	// Initialize Cloudinary
	log.Println("Initializing Cloudinary...")
	log.Printf("Cloud Name: %s", cfg.CloudinaryCloudName)
	log.Printf("API Key: %s", cfg.CloudinaryAPIKey)
	log.Printf("Upload Folder: %s", cfg.CloudinaryUploadFolder)
	
	if err := utils.InitCloudinary(); err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}
	log.Println("✓ Cloudinary initialized successfully")

	// Initialize Gin router
	log.Println("Initializing Gin router...")
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())
	log.Println("✓ CORS middleware added")

	// Setup routes
	routes.SetupRoutes(router, db)
	log.Println("✓ Routes configured")

	// Start server
	log.Printf("🚀 Server starting on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}