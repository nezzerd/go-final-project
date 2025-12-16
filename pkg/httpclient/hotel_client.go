package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HotelHTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewHotelHTTPClient(baseURL string) *HotelHTTPClient {
	return &HotelHTTPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Hotel struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
}

func (c *HotelHTTPClient) GetHotelOwnerID(ctx context.Context, hotelID string) (string, error) {
	url := fmt.Sprintf("%s/api/hotels/%s", c.baseURL, hotelID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get hotel: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("hotel service returned status %d: %s", resp.StatusCode, string(body))
	}

	var hotel Hotel
	if err := json.NewDecoder(resp.Body).Decode(&hotel); err != nil {
		return "", fmt.Errorf("failed to decode hotel: %w", err)
	}

	return hotel.OwnerID, nil
}
