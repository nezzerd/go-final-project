package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hotel-booking-system/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHotelHTTPClient(t *testing.T) {
	logger.Init("info")

	client := NewHotelHTTPClient("http://example.com")
	assert.NotNil(t, client)
	assert.Equal(t, "http://example.com", client.baseURL)
}

func TestHotelHTTPClient_GetHotelOwnerID(t *testing.T) {
	logger.Init("info")

	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/api/hotels/hotel-123", r.URL.Path)

			hotel := Hotel{
				ID:      "hotel-123",
				OwnerID: "owner-123",
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(hotel)
		}))
		defer server.Close()

		client := NewHotelHTTPClient(server.URL)

		ownerID, err := client.GetHotelOwnerID(context.Background(), "hotel-123")
		require.NoError(t, err)
		assert.Equal(t, "owner-123", ownerID)
	})

	t.Run("hotel not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("hotel not found"))
		}))
		defer server.Close()

		client := NewHotelHTTPClient(server.URL)

		ownerID, err := client.GetHotelOwnerID(context.Background(), "hotel-123")
		assert.Error(t, err)
		assert.Empty(t, ownerID)
		assert.Contains(t, err.Error(), "status 404")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client := NewHotelHTTPClient(server.URL)

		ownerID, err := client.GetHotelOwnerID(context.Background(), "hotel-123")
		assert.Error(t, err)
		assert.Empty(t, ownerID)
	})

	t.Run("invalid URL", func(t *testing.T) {
		client := NewHotelHTTPClient("http://invalid-url-that-does-not-exist:9999")

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		ownerID, err := client.GetHotelOwnerID(ctx, "hotel-123")
		assert.Error(t, err)
		assert.Empty(t, ownerID)
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}))
		defer server.Close()

		client := NewHotelHTTPClient(server.URL)

		ownerID, err := client.GetHotelOwnerID(context.Background(), "hotel-123")
		assert.Error(t, err)
		assert.Empty(t, ownerID)
		assert.Contains(t, err.Error(), "status 500")
	})
}

func TestHotel(t *testing.T) {
	hotel := Hotel{
		ID:      "hotel-123",
		OwnerID: "owner-123",
	}

	data, err := json.Marshal(hotel)
	require.NoError(t, err)

	var unmarshaled Hotel
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, hotel.ID, unmarshaled.ID)
	assert.Equal(t, hotel.OwnerID, unmarshaled.OwnerID)
}
