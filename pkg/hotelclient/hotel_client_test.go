package hotelclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHotelClient(t *testing.T) {
	client, err := NewHotelClient("localhost:8081")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	client.Close()
}

func TestHotelClient_GetRoomPrice(t *testing.T) {
	client, err := NewHotelClient("localhost:8081")
	assert.NoError(t, err)

	_, err = client.GetRoomPrice(context.Background(), "hotel-id", "room-id")
	assert.Error(t, err)

	client.Close()
}

func TestHotelClient_Close(t *testing.T) {
	client, err := NewHotelClient("localhost:8081")
	assert.NoError(t, err)

	err = client.Close()
	assert.NoError(t, err)
}
