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

func TestNewDeliveryClient(t *testing.T) {
	logger.Init("info")

	client := NewDeliveryClient("http://example.com")
	assert.NotNil(t, client)
	assert.Equal(t, "http://example.com", client.baseURL)
}

func TestDeliveryClient_SendNotification(t *testing.T) {
	logger.Init("info")

	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/notifications/send", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var req SendNotificationRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, "email", req.Channel)
			assert.Equal(t, "test@example.com", req.Recipient)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
			})
		}))
		defer server.Close()

		client := NewDeliveryClient(server.URL)

		req := &SendNotificationRequest{
			Channel:   "email",
			Recipient: "test@example.com",
			Subject:   "Test Subject",
			Message:   "test message",
		}

		err := client.SendNotification(context.Background(), req)
		assert.NoError(t, err)
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewDeliveryClient(server.URL)

		req := &SendNotificationRequest{
			Channel:   "email",
			Recipient: "test@example.com",
			Message:   "test message",
		}

		err := client.SendNotification(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status 500")
	})

	t.Run("invalid URL", func(t *testing.T) {
		client := NewDeliveryClient("http://invalid-url-that-does-not-exist:9999")

		req := &SendNotificationRequest{
			Channel:   "email",
			Recipient: "test@example.com",
			Message:   "test message",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		err := client.SendNotification(ctx, req)
		assert.Error(t, err)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		client := NewDeliveryClient("http://example.com")

		req := &SendNotificationRequest{
			Channel:   "email",
			Recipient: "test@example.com",
			Message:   "test message",
		}

		client.baseURL = "invalid-url"
		err := client.SendNotification(context.Background(), req)
		assert.Error(t, err)
	})
}

func TestSendNotificationRequest(t *testing.T) {
	req := &SendNotificationRequest{
		Channel:   "email",
		Recipient: "test@example.com",
		Subject:   "Test Subject",
		Message:   "test message",
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var unmarshaled SendNotificationRequest
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, req.Channel, unmarshaled.Channel)
	assert.Equal(t, req.Recipient, unmarshaled.Recipient)
	assert.Equal(t, req.Subject, unmarshaled.Subject)
	assert.Equal(t, req.Message, unmarshaled.Message)
}
