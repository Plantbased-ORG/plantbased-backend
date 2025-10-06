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
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Cloudinary
	if err := utils.InitCloudinary(); err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Setup routes
	routes.SetupRoutes(router, db)

	// Start server
	log.Printf("Server starting on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}