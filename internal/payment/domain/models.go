package domain

type PaymentRequest struct {
	BookingID string  `json:"booking_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency,omitempty"`
}

type PaymentResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
}

type PaymentWebhook struct {
	PaymentID   string  `json:"payment_id"`
	BookingID   string  `json:"booking_id"`
	Status      string  `json:"status"`
	Amount      float64 `json:"amount"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}
