package models

import "time"

// Admin represents an admin user
type Admin struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose in JSON
	FullName     string    `json:"full_name"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Admin        Admin  `json:"admin"`
}

// RefreshTokenRequest represents the refresh token payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UpdateProfileRequest represents profile update payload
type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email" binding:"email"`
}