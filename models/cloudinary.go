package models

// CloudinaryUploadResponse represents the response from Cloudinary upload
type CloudinaryUploadResponse struct {
	PublicID  string `json:"public_id"`
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
}