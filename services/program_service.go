package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"mime/multipart"
	"plantbased-backend/models"
	"plantbased-backend/utils"
)

type ProgramService struct {
	DB *sql.DB
}

func NewProgramService(db *sql.DB) *ProgramService {
	return &ProgramService{DB: db}
}

// CreateProgram creates a new program with images and pricing plans
func (s *ProgramService) CreateProgram(
	req models.CreateProgramRequest,
	images map[string]multipart.File,
) (*models.ProgramResponse, error) {

	// Upload all images to Cloudinary
	mainImage, err := utils.UploadImage(images["mainImage"], "programs")
	if err != nil {
		return nil, err
	}

	mainContentImage, err := utils.UploadImage(images["mainContentImage"], "programs")
	if err != nil {
		utils.DeleteImage(mainImage.PublicID)
		return nil, err
	}

	whatCausesImage, err := utils.UploadImage(images["whatCausesImage"], "programs")
	if err != nil {
		utils.DeleteImage(mainImage.PublicID)
		utils.DeleteImage(mainContentImage.PublicID)
		return nil, err
	}

	healthRisksImage, err := utils.UploadImage(images["healthRisksImage"], "programs")
	if err != nil {
		utils.DeleteImage(mainImage.PublicID)
		utils.DeleteImage(mainContentImage.PublicID)
		utils.DeleteImage(whatCausesImage.PublicID)
		return nil, err
	}

	strategiesImage, err := utils.UploadImage(images["strategiesImage"], "programs")
	if err != nil {
		utils.DeleteImage(mainImage.PublicID)
		utils.DeleteImage(mainContentImage.PublicID)
		utils.DeleteImage(whatCausesImage.PublicID)
		utils.DeleteImage(healthRisksImage.PublicID)
		return nil, err
	}

	conclusionImage, err := utils.UploadImage(images["conclusionImage"], "programs")
	if err != nil {
		utils.DeleteImage(mainImage.PublicID)
		utils.DeleteImage(mainContentImage.PublicID)
		utils.DeleteImage(whatCausesImage.PublicID)
		utils.DeleteImage(healthRisksImage.PublicID)
		utils.DeleteImage(strategiesImage.PublicID)
		return nil, err
	}

	// Insert program into database
	var programID int
	err = s.DB.QueryRow(`
		INSERT INTO programs (
			name, short_description, main_image_public_id, main_image_url,
			intro_description, main_content_image_public_id, main_content_image_url,
			what_causes, what_causes_image_public_id, what_causes_image_url,
			health_risks, health_risks_image_public_id, health_risks_image_url,
			strategies, strategies_image_public_id, strategies_image_url,
			conclusion, conclusion_image_public_id, conclusion_image_url
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id
	`, req.Name, req.ShortDescription, mainImage.PublicID, mainImage.SecureURL,
		req.IntroDescription, mainContentImage.PublicID, mainContentImage.SecureURL,
		req.WhatCauses, whatCausesImage.PublicID, whatCausesImage.SecureURL,
		req.HealthRisks, healthRisksImage.PublicID, healthRisksImage.SecureURL,
		req.Strategies, strategiesImage.PublicID, strategiesImage.SecureURL,
		req.Conclusion, conclusionImage.PublicID, conclusionImage.SecureURL,
	).Scan(&programID)

	if err != nil {
		// Rollback: delete all uploaded images
		utils.DeleteImage(mainImage.PublicID)
		utils.DeleteImage(mainContentImage.PublicID)
		utils.DeleteImage(whatCausesImage.PublicID)
		utils.DeleteImage(healthRisksImage.PublicID)
		utils.DeleteImage(strategiesImage.PublicID)
		utils.DeleteImage(conclusionImage.PublicID)
		return nil, err
	}

	// Insert pricing plans
	var pricingPlans []models.ProgramPricingPlan
	for _, plan := range req.PricingPlans {
		featuresJSON, _ := json.Marshal(plan.Features)

		var planID int
		err = s.DB.QueryRow(`
			INSERT INTO program_pricing_plans (program_id, name, subtitle, price, features)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, programID, plan.Name, plan.Subtitle, plan.Price, featuresJSON).Scan(&planID)

		if err != nil {
			return nil, err
		}

		pricingPlans = append(pricingPlans, models.ProgramPricingPlan{
			ID:        planID,
			ProgramID: programID,
			Name:      plan.Name,
			Subtitle:  plan.Subtitle,
			Price:     plan.Price,
			Features:  plan.Features,
		})
	}

	// Get the created program
	program, err := s.GetProgramByID(programID)
	if err != nil {
		return nil, err
	}

	return &models.ProgramResponse{
		Program:      *program,
		PricingPlans: pricingPlans,
	}, nil
}

// GetAllPrograms retrieves all programs with their pricing plans
func (s *ProgramService) GetAllPrograms() ([]models.ProgramResponse, error) {
	rows, err := s.DB.Query(`
		SELECT id, name, short_description, main_image_public_id, main_image_url,
		intro_description, main_content_image_public_id, main_content_image_url,
		what_causes, what_causes_image_public_id, what_causes_image_url,
		health_risks, health_risks_image_public_id, health_risks_image_url,
		strategies, strategies_image_public_id, strategies_image_url,
		conclusion, conclusion_image_public_id, conclusion_image_url,
		created_at, updated_at
		FROM programs
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []models.ProgramResponse
	for rows.Next() {
		var p models.Program
		err := rows.Scan(
			&p.ID, &p.Name, &p.ShortDescription, &p.MainImagePublicID, &p.MainImageURL,
			&p.IntroDescription, &p.MainContentImagePublicID, &p.MainContentImageURL,
			&p.WhatCauses, &p.WhatCausesImagePublicID, &p.WhatCausesImageURL,
			&p.HealthRisks, &p.HealthRisksImagePublicID, &p.HealthRisksImageURL,
			&p.Strategies, &p.StrategiesImagePublicID, &p.StrategiesImageURL,
			&p.Conclusion, &p.ConclusionImagePublicID, &p.ConclusionImageURL,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get pricing plans for this program
		pricingPlans, err := s.GetPricingPlansByProgramID(p.ID)
		if err != nil {
			return nil, err
		}

		programs = append(programs, models.ProgramResponse{
			Program:      p,
			PricingPlans: pricingPlans,
		})
	}

	return programs, nil
}

// GetProgramByID retrieves a single program by ID
func (s *ProgramService) GetProgramByID(id int) (*models.Program, error) {
	var p models.Program
	err := s.DB.QueryRow(`
		SELECT id, name, short_description, main_image_public_id, main_image_url,
		intro_description, main_content_image_public_id, main_content_image_url,
		what_causes, what_causes_image_public_id, what_causes_image_url,
		health_risks, health_risks_image_public_id, health_risks_image_url,
		strategies, strategies_image_public_id, strategies_image_url,
		conclusion, conclusion_image_public_id, conclusion_image_url,
		created_at, updated_at
		FROM programs WHERE id = $1
	`, id).Scan(
		&p.ID, &p.Name, &p.ShortDescription, &p.MainImagePublicID, &p.MainImageURL,
		&p.IntroDescription, &p.MainContentImagePublicID, &p.MainContentImageURL,
		&p.WhatCauses, &p.WhatCausesImagePublicID, &p.WhatCausesImageURL,
		&p.HealthRisks, &p.HealthRisksImagePublicID, &p.HealthRisksImageURL,
		&p.Strategies, &p.StrategiesImagePublicID, &p.StrategiesImageURL,
		&p.Conclusion, &p.ConclusionImagePublicID, &p.ConclusionImageURL,
		&p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("program not found")
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetPricingPlansByProgramID retrieves all pricing plans for a program
func (s *ProgramService) GetPricingPlansByProgramID(programID int) ([]models.ProgramPricingPlan, error) {
	rows, err := s.DB.Query(`
		SELECT id, program_id, name, subtitle, price, features, created_at, updated_at
		FROM program_pricing_plans
		WHERE program_id = $1
		ORDER BY id ASC
	`, programID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.ProgramPricingPlan
	for rows.Next() {
		var plan models.ProgramPricingPlan
		var featuresJSON []byte

		err := rows.Scan(
			&plan.ID, &plan.ProgramID, &plan.Name, &plan.Subtitle,
			&plan.Price, &featuresJSON, &plan.CreatedAt, &plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal features JSON
		json.Unmarshal(featuresJSON, &plan.Features)
		plans = append(plans, plan)
	}

	return plans, nil
}

// DeleteProgram deletes a program and its images
func (s *ProgramService) DeleteProgram(id int) error {
	// Get program to retrieve image public IDs
	program, err := s.GetProgramByID(id)
	if err != nil {
		return err
	}

	// Delete from database (pricing plans will be cascade deleted)
	_, err = s.DB.Exec("DELETE FROM programs WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Delete images from Cloudinary
	utils.DeleteImage(program.MainImagePublicID)
	utils.DeleteImage(program.MainContentImagePublicID)
	utils.DeleteImage(program.WhatCausesImagePublicID)
	utils.DeleteImage(program.HealthRisksImagePublicID)
	utils.DeleteImage(program.StrategiesImagePublicID)
	utils.DeleteImage(program.ConclusionImagePublicID)

	return nil
}