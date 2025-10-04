package utils

import (
	"plantbased-backend/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a JWT token
func GenerateToken(adminID int, email string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": adminID,
		"email":    email,
		"exp":      time.Now().Add(expiry).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})
}