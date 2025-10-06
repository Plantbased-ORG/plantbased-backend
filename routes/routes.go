package routes

import (
	"database/sql"
	"plantbased-backend/handlers"
	"plantbased-backend/middleware"
	"plantbased-backend/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// Initialize services
	authService := services.NewAuthService(db)
	adminService := services.NewAdminService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Admin routes (protected)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		{
			admin.GET("/profile", adminHandler.GetProfile)
			admin.PUT("/profile", adminHandler.UpdateProfile)
		}
	}
}