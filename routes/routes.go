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
	programService := services.NewProgramService(db)
	testimonialService := services.NewTestimonialService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)
	programHandler := handlers.NewProgramHandler(programService)
	testimonialHandler := handlers.NewTestimonialHandler(testimonialService)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes (public)
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

		// Program routes (protected - only admin can create/delete)
		programs := api.Group("/programs")
		{
			programs.GET("", programHandler.GetAllPrograms)           // Public
			programs.GET("/:id", programHandler.GetProgramByID)       // Public
			programs.POST("", middleware.AuthMiddleware(), programHandler.CreateProgram)      // Protected
			programs.DELETE("/:id", middleware.AuthMiddleware(), programHandler.DeleteProgram) // Protected
		}

		// Testimonial routes (protected - only admin can create/update/delete)
		testimonials := api.Group("/testimonials")
		{
			testimonials.GET("", testimonialHandler.GetAllTestimonials)           // Public
			testimonials.GET("/:id", testimonialHandler.GetTestimonialByID)       // Public
			testimonials.POST("", middleware.AuthMiddleware(), testimonialHandler.CreateTestimonial)      // Protected
			testimonials.PUT("/:id", middleware.AuthMiddleware(), testimonialHandler.UpdateTestimonial)   // Protected
			testimonials.DELETE("/:id", middleware.AuthMiddleware(), testimonialHandler.DeleteTestimonial) // Protected
		}
	}
}