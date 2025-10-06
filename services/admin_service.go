package services

import (
	"database/sql"
	"errors"
	"plantbased-backend/models"
	"plantbased-backend/utils"
)

type AdminService struct {
	db *sql.DB
}

func NewAdminService(db *sql.DB) *AdminService {
	return &AdminService{db: db}
}

func (s *AdminService) GetAdminByID(id int) (*models.Admin, error) {
	var admin models.Admin
	query := `SELECT id, email, name, created_at, updated_at FROM admins WHERE id = $1`
	
	err := s.db.QueryRow(query, id).Scan(
		&admin.ID,
		&admin.Email,
		&admin.Name,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("admin not found")
	}
	if err != nil {
		return nil, err
	}
	
	return &admin, nil
}

func (s *AdminService) UpdateAdmin(id int, req models.UpdateAdminRequest) (*models.Admin, error) {
	// Build dynamic update query
	query := `UPDATE admins SET name = $1, updated_at = NOW() WHERE id = $2 RETURNING id, email, name, created_at, updated_at`
	
	var admin models.Admin
	err := s.db.QueryRow(query, req.Name, id).Scan(
		&admin.ID,
		&admin.Email,
		&admin.Name,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &admin, nil
}

func (s *AdminService) UpdatePassword(id int, currentPassword, newPassword string) error {
	// Get current password hash
	var passwordHash string
	query := `SELECT password_hash FROM admins WHERE id = $1`
	err := s.db.QueryRow(query, id).Scan(&passwordHash)
	if err != nil {
		return errors.New("admin not found")
	}
	
	// Verify current password
	if !utils.CheckPasswordHash(currentPassword, passwordHash) {
		return errors.New("incorrect current password")
	}
	
	// Hash new password
	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	
	// Update password
	updateQuery := `UPDATE admins SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err = s.db.Exec(updateQuery, newHash, id)
	return err
}