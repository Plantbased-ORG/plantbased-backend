package models

// CustomerDetails represents customer registration data
type CustomerDetails struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	Nationality string `json:"nationality"`
	PhoneNumber string `json:"phoneNumber"`
	Program     string `json:"program"`
	Package     string `json:"package"`
}