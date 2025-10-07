package handlers

import (
	"net/http"
	"plantbased-backend/models"
	"plantbased-backend/services"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) GetProfile(c *gin.Context) {
	adminID := c.GetInt("adminID")
	
	admin, err := h.adminService.GetAdminByID(adminID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Admin not found",
		})
		return
	}
	
	c.JSON(http.StatusOK, admin)
}

func (h *AdminHandler) UpdateProfile(c *gin.Context) {
	adminID := c.GetInt("adminID")
	
	var req models.UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	
	admin, err := h.adminService.UpdateAdmin(adminID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update profile",
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, admin)
}

func (h *AdminHandler) ChangePassword(c *gin.Context) {
	adminID := c.GetInt("adminID")
	
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	
	// Validate that new password is different from current
	if req.CurrentPassword == req.NewPassword {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "New password must be different from current password",
		})
		return
	}
	
	err := h.adminService.UpdatePassword(adminID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to change password",
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}