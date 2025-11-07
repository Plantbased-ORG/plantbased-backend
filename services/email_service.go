package services

import (
	"fmt"
	"net/smtp"
	"os"
	"plantbased-backend/models"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendCustomerDetailsToCEO(details models.CustomerDetails) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	ceoEmail := os.Getenv("CEO_EMAIL")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := fmt.Sprintf("New Customer Registration: %s", details.FullName)
	body := fmt.Sprintf(`New customer has registered for PlantBased Meals:

Full Name: %s
Email: %s
Phone Number: %s
Nationality: %s
Program: %s
Package: %s

Best regards,
PlantBased Meals System`,
		details.FullName,
		details.Email,
		details.PhoneNumber,
		details.Nationality,
		details.Program,
		details.Package,
	)

	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{ceoEmail},
		message,
	)

	return err
}