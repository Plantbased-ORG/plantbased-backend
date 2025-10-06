package models

import "time"

// Program represents a healing program
type Program struct {
	ID                      int       `json:"id"`
	Name                    string    `json:"name"`
	ShortDescription        string    `json:"short_description"`
	MainImagePublicID       string    `json:"main_image_public_id"`
	MainImageURL            string    `json:"main_image_url"`
	IntroDescription        string    `json:"intro_description"`
	MainContentImagePublicID string   `json:"main_content_image_public_id"`
	MainContentImageURL     string    `json:"main_content_image_url"`
	WhatCauses              string    `json:"what_causes"`
	WhatCausesImagePublicID string    `json:"what_causes_image_public_id"`
	WhatCausesImageURL      string    `json:"what_causes_image_url"`
	HealthRisks             string    `json:"health_risks"`
	HealthRisksImagePublicID string   `json:"health_risks_image_public_id"`
	HealthRisksImageURL     string    `json:"health_risks_image_url"`
	Strategies              string    `json:"strategies"`
	StrategiesImagePublicID string    `json:"strategies_image_public_id"`
	StrategiesImageURL      string    `json:"strategies_image_url"`
	Conclusion              string    `json:"conclusion"`
	ConclusionImagePublicID string    `json:"conclusion_image_public_id"`
	ConclusionImageURL      string    `json:"conclusion_image_url"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// ProgramPricingPlan represents a pricing plan for a program
type ProgramPricingPlan struct {
	ID        int       `json:"id"`
	ProgramID int       `json:"program_id"`
	Name      string    `json:"name"`
	Subtitle  string    `json:"subtitle"`
	Price     string    `json:"price"`
	Features  []string  `json:"features"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateProgramRequest represents the request to create a program
type CreateProgramRequest struct {
	Name             string               `json:"name"`
	ShortDescription string               `json:"short_description"`
	IntroDescription string               `json:"intro_description"`
	WhatCauses       string               `json:"what_causes"`
	HealthRisks      string               `json:"health_risks"`
	Strategies       string               `json:"strategies"`
	Conclusion       string               `json:"conclusion"`
	PricingPlans     []PricingPlanRequest `json:"pricing_plans"`
}

// PricingPlanRequest represents a pricing plan in the request
type PricingPlanRequest struct {
	Name     string   `json:"name"`
	Subtitle string   `json:"subtitle"`
	Price    string   `json:"price"`
	Features []string `json:"features"`
}

// ProgramResponse represents a complete program with pricing plans
type ProgramResponse struct {
	Program      Program              `json:"program"`
	PricingPlans []ProgramPricingPlan `json:"pricing_plans"`
}