package handlers

import (
	"plantbased-backend/models"
	"plantbased-backend/services"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	emailService *services.EmailService
}

func NewCustomerHandler(emailService *services.EmailService) *CustomerHandler {
	return &CustomerHandler{emailService: emailService}
}

func (h *CustomerHandler) SendCustomerDetails(c *gin.Context) {
	var details models.CustomerDetails
	
	if err := c.ShouldBindJSON(&details); err != nil {
		c.JSON(400, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}

	err := h.emailService.SendCustomerDetailsToCEO(details)
	if err != nil {
		c.JSON(500, models.ErrorResponse{
			Error:   "email_failed",
			Message: "Failed to send email",
		})
		return
	}

	c.JSON(200, models.SuccessResponse{
		Success: true,
		Message: "Customer details sent successfully",
	})
}