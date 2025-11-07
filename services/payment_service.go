package services

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"os"
)

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// VerifyWebhookSignature verifies that the webhook came from Paystack
func (s *PaymentService) VerifyWebhookSignature(signature string, payload []byte) bool {
	secret := os.Getenv("PAYSTACK_SECRET_KEY")
	
	hash := hmac.New(sha512.New, []byte(secret))
	hash.Write(payload)
	expectedSignature := hex.EncodeToString(hash.Sum(nil))
	
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}