package handlers

import (
	"net/http"
	"plantbased-backend/models"
	"plantbased-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TestimonialHandler struct {
	testimonialService *services.TestimonialService
}

func NewTestimonialHandler(testimonialService *services.TestimonialService) *TestimonialHandler {
	return &TestimonialHandler{testimonialService: testimonialService}
}

// CreateTestimonial handles testimonial creation
func (h *TestimonialHandler) CreateTestimonial(c *gin.Context) {
	var req models.CreateTestimonialRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	testimonial, err := h.testimonialService.CreateTestimonial(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create testimonial",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, testimonial)
}

// GetAllTestimonials retrieves all testimonials
func (h *TestimonialHandler) GetAllTestimonials(c *gin.Context) {
	testimonials, err := h.testimonialService.GetAllTestimonials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch testimonials",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, testimonials)
}

// GetTestimonialByID retrieves a single testimonial
func (h *TestimonialHandler) GetTestimonialByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid testimonial ID",
		})
		return
	}

	testimonial, err := h.testimonialService.GetTestimonialByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, testimonial)
}

// UpdateTestimonial updates an existing testimonial
func (h *TestimonialHandler) UpdateTestimonial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid testimonial ID",
		})
		return
	}

	var req models.CreateTestimonialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	testimonial, err := h.testimonialService.UpdateTestimonial(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update testimonial",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Testimonial updated successfully",
		Data:    testimonial,
	})
}

// DeleteTestimonial deletes a testimonial
func (h *TestimonialHandler) DeleteTestimonial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid testimonial ID",
		})
		return
	}

	err = h.testimonialService.DeleteTestimonial(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete testimonial",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Testimonial deleted successfully",
	})
}