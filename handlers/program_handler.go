package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"plantbased-backend/models"
	"plantbased-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProgramHandler struct {
	programService *services.ProgramService
}

func NewProgramHandler(programService *services.ProgramService) *ProgramHandler {
	return &ProgramHandler{programService: programService}
}

// CreateProgram handles program creation with images
func (h *ProgramHandler) CreateProgram(c *gin.Context) {
	// Parse multipart form
	err := c.Request.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid form data",
			Message: err.Error(),
		})
		return
	}

	// Get text fields
	name := c.PostForm("name")
	shortDescription := c.PostForm("shortDescription")
	introDescription := c.PostForm("introDescription")
	whatCauses := c.PostForm("whatCauses")
	healthRisks := c.PostForm("healthRisks")
	strategies := c.PostForm("strategies")
	conclusion := c.PostForm("conclusion")
	pricingPlansJSON := c.PostForm("pricingPlans")

	// Validate required fields
	if name == "" || shortDescription == "" || pricingPlansJSON == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing required fields",
		})
		return
	}

	// Parse pricing plans
	var pricingPlans []models.PricingPlanRequest
	if err := json.Unmarshal([]byte(pricingPlansJSON), &pricingPlans); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid pricing plans format",
			Message: err.Error(),
		})
		return
	}

	// Get image files
	images := make(map[string]multipart.File)
	imageFields := []string{
		"mainImage", "mainContentImage", "whatCausesImage",
		"healthRisksImage", "strategiesImage", "conclusionImage",
	}

	for _, field := range imageFields {
		file, _, err := c.Request.FormFile(field)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Missing required image: " + field,
				Message: err.Error(),
			})
			return
		}
		images[field] = file
		defer file.Close()
	}

	// Create request
	req := models.CreateProgramRequest{
		Name:             name,
		ShortDescription: shortDescription,
		IntroDescription: introDescription,
		WhatCauses:       whatCauses,
		HealthRisks:      healthRisks,
		Strategies:       strategies,
		Conclusion:       conclusion,
		PricingPlans:     pricingPlans,
	}

	// Create program
	response, err := h.programService.CreateProgram(req, images)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create program",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateProgram handles program updates with optional image changes
func (h *ProgramHandler) UpdateProgram(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid program ID",
		})
		return
	}

	// Parse multipart form
	err = c.Request.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid form data",
			Message: err.Error(),
		})
		return
	}

	// Get text fields
	name := c.PostForm("name")
	shortDescription := c.PostForm("shortDescription")
	introDescription := c.PostForm("introDescription")
	whatCauses := c.PostForm("whatCauses")
	healthRisks := c.PostForm("healthRisks")
	strategies := c.PostForm("strategies")
	conclusion := c.PostForm("conclusion")
	pricingPlansJSON := c.PostForm("pricingPlans")

	// Validate required fields
	if name == "" || shortDescription == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing required fields",
		})
		return
	}

	// Parse pricing plans (optional)
	var pricingPlans []models.PricingPlanRequest
	if pricingPlansJSON != "" {
		if err := json.Unmarshal([]byte(pricingPlansJSON), &pricingPlans); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid pricing plans format",
				Message: err.Error(),
			})
			return
		}
	}

	// Get image files (all optional for updates)
	images := make(map[string]multipart.File)
	imageFields := []string{
		"mainImage", "mainContentImage", "whatCausesImage",
		"healthRisksImage", "strategiesImage", "conclusionImage",
	}

	for _, field := range imageFields {
		file, _, err := c.Request.FormFile(field)
		if err == nil {
			images[field] = file
			defer file.Close()
		}
	}

	// Create request
	req := models.CreateProgramRequest{
		Name:             name,
		ShortDescription: shortDescription,
		IntroDescription: introDescription,
		WhatCauses:       whatCauses,
		HealthRisks:      healthRisks,
		Strategies:       strategies,
		Conclusion:       conclusion,
		PricingPlans:     pricingPlans,
	}

	// Update program
	response, err := h.programService.UpdateProgram(id, req, images)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update program",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAllPrograms retrieves all programs
func (h *ProgramHandler) GetAllPrograms(c *gin.Context) {
	programs, err := h.programService.GetAllPrograms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch programs",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, programs)
}

// GetProgramByID retrieves a single program
func (h *ProgramHandler) GetProgramByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid program ID",
		})
		return
	}

	program, err := h.programService.GetProgramByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Get pricing plans
	pricingPlans, err := h.programService.GetPricingPlansByProgramID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch pricing plans",
		})
		return
	}

	c.JSON(http.StatusOK, models.ProgramResponse{
		Program:      *program,
		PricingPlans: pricingPlans,
	})
}

// DeleteProgram deletes a program
func (h *ProgramHandler) DeleteProgram(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid program ID",
		})
		return
	}

	err = h.programService.DeleteProgram(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete program",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Program deleted successfully",
	})
}

// AddPricingPlan adds a pricing plan to an existing program
func (h *ProgramHandler) AddPricingPlan(c *gin.Context) {
	programID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid program ID",
		})
		return
	}

	var req models.PricingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	plan, err := h.programService.AddPricingPlan(programID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to add pricing plan",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

// UpdatePricingPlan updates a specific pricing plan
func (h *ProgramHandler) UpdatePricingPlan(c *gin.Context) {
	programID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid program ID",
		})
		return
	}

	planID, err := strconv.Atoi(c.Param("plan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid pricing plan ID",
		})
		return
	}

	var req models.PricingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	plan, err := h.programService.UpdatePricingPlan(programID, planID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update pricing plan",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// DeletePricingPlan deletes a specific pricing plan
func (h *ProgramHandler) DeletePricingPlan(c *gin.Context) {
	programID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid program ID",
		})
		return
	}

	planID, err := strconv.Atoi(c.Param("plan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid pricing plan ID",
		})
		return
	}

	err = h.programService.DeletePricingPlan(programID, planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete pricing plan",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Pricing plan deleted successfully",
	})
}