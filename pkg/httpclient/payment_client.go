package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"hotel-booking-system/pkg/logger"
)

type PaymentClient struct {
	baseURL string
	client  *http.Client
}

func NewPaymentClient(baseURL string) *PaymentClient {
	return &PaymentClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

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

func (c *PaymentClient) CreatePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/api/payments", c.baseURL)

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to create payment")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("payment service returned status %d", resp.StatusCode)
	}

	var paymentResp PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &paymentResp, nil
}
