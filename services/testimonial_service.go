package services

import (
	"database/sql"
	"errors"
	"plantbased-backend/models"
)

type TestimonialService struct {
	DB *sql.DB
}

func NewTestimonialService(db *sql.DB) *TestimonialService {
	return &TestimonialService{DB: db}
}

// CreateTestimonial creates a new testimonial
func (s *TestimonialService) CreateTestimonial(req models.CreateTestimonialRequest) (*models.Testimonial, error) {
	var testimonial models.Testimonial

	err := s.DB.QueryRow(`
		INSERT INTO testimonials (name, location, review, avatar)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, location, review, avatar, created_at, updated_at
	`, req.Name, req.Location, req.Review, req.Avatar).Scan(
		&testimonial.ID,
		&testimonial.Name,
		&testimonial.Location,
		&testimonial.Review,
		&testimonial.Avatar,
		&testimonial.CreatedAt,
		&testimonial.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &testimonial, nil
}

// GetAllTestimonials retrieves all testimonials
func (s *TestimonialService) GetAllTestimonials() ([]models.Testimonial, error) {
	rows, err := s.DB.Query(`
		SELECT id, name, location, review, avatar, created_at, updated_at
		FROM testimonials
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testimonials []models.Testimonial
	for rows.Next() {
		var t models.Testimonial
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Location,
			&t.Review,
			&t.Avatar,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		testimonials = append(testimonials, t)
	}

	return testimonials, nil
}

// GetTestimonialByID retrieves a single testimonial by ID
func (s *TestimonialService) GetTestimonialByID(id int) (*models.Testimonial, error) {
	var t models.Testimonial

	err := s.DB.QueryRow(`
		SELECT id, name, location, review, avatar, created_at, updated_at
		FROM testimonials
		WHERE id = $1
	`, id).Scan(
		&t.ID,
		&t.Name,
		&t.Location,
		&t.Review,
		&t.Avatar,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("testimonial not found")
	}

	if err != nil {
		return nil, err
	}

	return &t, nil
}

// UpdateTestimonial updates an existing testimonial
func (s *TestimonialService) UpdateTestimonial(id int, req models.CreateTestimonialRequest) (*models.Testimonial, error) {
	_, err := s.DB.Exec(`
		UPDATE testimonials
		SET name = $1, location = $2, review = $3, avatar = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`, req.Name, req.Location, req.Review, req.Avatar, id)

	if err != nil {
		return nil, err
	}

	return s.GetTestimonialByID(id)
}

// DeleteTestimonial deletes a testimonial
func (s *TestimonialService) DeleteTestimonial(id int) error {
	result, err := s.DB.Exec("DELETE FROM testimonials WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("testimonial not found")
	}

	return nil
}