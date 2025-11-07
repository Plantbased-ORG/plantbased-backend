package handlers

import (
	"encoding/json"
	"net/http"
	"plantbased-backend/models"
	"plantbased-backend/services"
)

type CustomerHandler struct {
	emailService *services.EmailService
}

func NewCustomerHandler(emailService *services.EmailService) *CustomerHandler {
	return &CustomerHandler{emailService: emailService}
}

func (h *CustomerHandler) SendCustomerDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var details models.CustomerDetails
	err := json.NewDecoder(r.Body).Decode(&details)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}

	err = h.emailService.SendCustomerDetailsToCEO(details)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "email_failed",
			Message: "Failed to send email",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Message: "Customer details sent successfully",
	})
}