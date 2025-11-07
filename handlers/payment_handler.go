package handlers

import (
	"io"
	"plantbased-backend/models"
	"plantbased-backend/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	signature := c.GetHeader("x-paystack-signature")
	if signature == "" {
		c.JSON(400, models.ErrorResponse{
			Error:   "invalid_signature",
			Message: "Missing signature",
		})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Failed to read request body",
		})
		return
	}

	if !h.paymentService.VerifyWebhookSignature(signature, body) {
		c.JSON(401, models.ErrorResponse{
			Error:   "invalid_signature",
			Message: "Invalid webhook signature",
		})
		return
	}

	var webhook models.PaystackWebhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(400, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid webhook data",
		})
		return
	}

	if webhook.Event == "charge.success" && webhook.Data.Status == "success" {
		// Payment successful - handle it here
		// TODO: Store order in database, send confirmation email, etc.
	}

	c.JSON(200, gin.H{"status": "success"})
}