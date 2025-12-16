package hotelclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HotelClient struct {
	baseURL string
}

func NewHotelClient(addr string) (*HotelClient, error) {
	return &HotelClient{
		baseURL: fmt.Sprintf("http://%s", addr),
	}, nil
}

func (c *HotelClient) GetRoomPrice(ctx context.Context, hotelID, roomID string) (float64, error) {
	url := fmt.Sprintf("%s/api/hotels/%s/rooms", c.baseURL, hotelID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("hotel service returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Hotel struct{} `json:"hotel"`
		Rooms []struct {
			ID            string  `json:"id"`
			PricePerNight float64 `json:"price_per_night"`
		} `json:"rooms"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse hotel service response: %w", err)
	}

	for _, room := range result.Rooms {
		if room.ID == roomID {
			return room.PricePerNight, nil
		}
	}

	return 0, fmt.Errorf("room %s not found in hotel %s", roomID, hotelID)
}

func (c *HotelClient) Close() error {
	return nil
}
