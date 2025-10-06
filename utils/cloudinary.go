package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"plantbased-backend/config"
	"plantbased-backend/models"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

// InitCloudinary initializes the Cloudinary client
func InitCloudinary() error {
	var err error
	cld, err = cloudinary.NewFromParams(
		config.AppConfig.CloudinaryCloudName,
		config.AppConfig.CloudinaryAPIKey,
		config.AppConfig.CloudinaryAPISecret,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}
	return nil
}

// UploadImage uploads an image to Cloudinary
func UploadImage(file multipart.File, folder string) (*models.CloudinaryUploadResponse, error) {
	ctx := context.Background()

	// Upload the file to Cloudinary
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: fmt.Sprintf("%s/%s", config.AppConfig.CloudinaryUploadFolder, folder),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	return &models.CloudinaryUploadResponse{
		PublicID:  uploadResult.PublicID,
		URL:       uploadResult.URL,
		SecureURL: uploadResult.SecureURL,
	}, nil
}

// DeleteImage deletes an image from Cloudinary by public ID
func DeleteImage(publicID string) error {
	ctx := context.Background()

	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}