package services

import (
	"database/sql"
	"errors"
	"plantbased-backend/models"
	"plantbased-backend/utils"
	"time"
)

// AuthService handles authentication business logic
type AuthService struct {
	DB *sql.DB
}

// NewAuthService creates a new AuthService
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{DB: db}
}

// Login authenticates an admin and returns tokens
func (s *AuthService) Login(email, password string) (*models.LoginResponse, error) {
	// Find admin by email
	var admin models.Admin
	err := s.DB.QueryRow(`
		SELECT id, email, password_hash, full_name, is_active, created_at, updated_at
		FROM admins
		WHERE email = $1
	`, email).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&admin.FullName,
		&admin.IsActive,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid email or password")
	}

	if err != nil {
		return nil, err
	}

	// Check if admin is active
	if !admin.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Verify password
	if err := utils.ComparePassword(admin.PasswordHash, password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	token, err := utils.GenerateToken(admin.ID, admin.Email, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(admin.ID, admin.Email, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Admin:        admin,
	}, nil
}

// RefreshToken generates a new access token from refresh token
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	// Validate refresh token
	token, err := utils.ValidateToken(refreshToken)
	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired refresh token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	adminID := int(claims["admin_id"].(float64))
	email := claims["email"].(string)

	// Generate new access token
	newToken, err := utils.GenerateToken(adminID, email, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return newToken, nil
}