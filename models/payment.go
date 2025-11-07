package models

// PaystackWebhook represents the webhook payload from Paystack
type PaystackWebhook struct {
	Event string             `json:"event"`
	Data  PaystackWebhookData `json:"data"`
}

// PaystackWebhookData represents the data inside the webhook
type PaystackWebhookData struct {
	Reference string                 `json:"reference"`
	Amount    int                    `json:"amount"`
	Status    string                 `json:"status"`
	Customer  PaystackCustomer       `json:"customer"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// PaystackCustomer represents customer info from Paystack
type PaystackCustomer struct {
	Email string `json:"email"`
}