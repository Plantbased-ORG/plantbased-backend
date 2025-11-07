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
	emailService := services.NewEmailService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)
	programHandler := handlers.NewProgramHandler(programService)
	testimonialHandler := handlers.NewTestimonialHandler(testimonialService)
	customerHandler := handlers.NewCustomerHandler(emailService)

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
			admin.PUT("/change-password", adminHandler.ChangePassword)
		}

		// Program routes
		programs := api.Group("/programs")
		{
			// Public routes
			programs.GET("", programHandler.GetAllPrograms)
			programs.GET("/:id", programHandler.GetProgramByID)

			// Protected routes (admin only)
			programs.POST("", middleware.AuthMiddleware(), programHandler.CreateProgram)
			programs.PUT("/:id", middleware.AuthMiddleware(), programHandler.UpdateProgram)
			programs.DELETE("/:id", middleware.AuthMiddleware(), programHandler.DeleteProgram)

			// Pricing plan routes (admin only)
			programs.POST("/:id/pricing-plans", middleware.AuthMiddleware(), programHandler.AddPricingPlan)
			programs.PUT("/:id/pricing-plans/:plan_id", middleware.AuthMiddleware(), programHandler.UpdatePricingPlan)
			programs.DELETE("/:id/pricing-plans/:plan_id", middleware.AuthMiddleware(), programHandler.DeletePricingPlan)
		}

		// Testimonial routes
		testimonials := api.Group("/testimonials")
		{
			// Public routes
			testimonials.GET("", testimonialHandler.GetAllTestimonials)
			testimonials.GET("/:id", testimonialHandler.GetTestimonialByID)

			// Protected routes (admin only)
			testimonials.POST("", middleware.AuthMiddleware(), testimonialHandler.CreateTestimonial)
			testimonials.PUT("/:id", middleware.AuthMiddleware(), testimonialHandler.UpdateTestimonial)
			testimonials.DELETE("/:id", middleware.AuthMiddleware(), testimonialHandler.DeleteTestimonial)
		}

		// Customer routes (public)
		api.POST("/send-customer-details", customerHandler.SendCustomerDetails)
	}
}