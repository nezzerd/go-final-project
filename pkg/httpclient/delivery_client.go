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

type DeliveryClient struct {
	baseURL string
	client  *http.Client
}

func NewDeliveryClient(baseURL string) *DeliveryClient {
	return &DeliveryClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type SendNotificationRequest struct {
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
	Subject   string `json:"subject,omitempty"`
	Message   string `json:"message"`
}

func (c *DeliveryClient) SendNotification(ctx context.Context, req *SendNotificationRequest) error {
	url := fmt.Sprintf("%s/api/notifications/send", c.baseURL)

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to send notification via delivery service")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delivery service returned status %d", resp.StatusCode)
	}

	return nil
}
