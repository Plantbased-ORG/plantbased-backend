package handlers

import (
	"net/http"
	"plantbased-backend/models"
	"plantbased-backend/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles admin login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Authenticate admin
	response, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   err.Error(),
			Message: "Please check your credentials",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Refresh token
	newToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}