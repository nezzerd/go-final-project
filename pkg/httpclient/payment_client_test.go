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

func TestNewPaymentClient(t *testing.T) {
	logger.Init("info")

	client := NewPaymentClient("http://example.com")
	assert.NotNil(t, client)
	assert.Equal(t, "http://example.com", client.baseURL)
}

func TestPaymentClient_CreatePayment(t *testing.T) {
	logger.Init("info")

	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/payments", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var req PaymentRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.Equal(t, "booking-123", req.BookingID)
			assert.Equal(t, 1000.0, req.Amount)

			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(PaymentResponse{
				PaymentID: "payment-123",
				Status:    "processing",
				Message:   "payment is being processed",
			})
		}))
		defer server.Close()

		client := NewPaymentClient(server.URL)

		req := &PaymentRequest{
			BookingID: "booking-123",
			Amount:    1000.0,
			Currency:  "RUB",
		}

		response, err := client.CreatePayment(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "payment-123", response.PaymentID)
		assert.Equal(t, "processing", response.Status)
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewPaymentClient(server.URL)

		req := &PaymentRequest{
			BookingID: "booking-123",
			Amount:    1000.0,
		}

		response, err := client.CreatePayment(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "status 500")
	})

	t.Run("invalid status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewPaymentClient(server.URL)

		req := &PaymentRequest{
			BookingID: "booking-123",
			Amount:    1000.0,
		}

		response, err := client.CreatePayment(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "status 200")
	})

	t.Run("invalid URL", func(t *testing.T) {
		client := NewPaymentClient("http://invalid-url-that-does-not-exist:9999")

		req := &PaymentRequest{
			BookingID: "booking-123",
			Amount:    1000.0,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		response, err := client.CreatePayment(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client := NewPaymentClient(server.URL)

		req := &PaymentRequest{
			BookingID: "booking-123",
			Amount:    1000.0,
		}

		response, err := client.CreatePayment(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestPaymentRequest(t *testing.T) {
	req := &PaymentRequest{
		BookingID: "booking-123",
		Amount:    1000.0,
		Currency:  "RUB",
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var unmarshaled PaymentRequest
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, req.BookingID, unmarshaled.BookingID)
	assert.Equal(t, req.Amount, unmarshaled.Amount)
	assert.Equal(t, req.Currency, unmarshaled.Currency)
}
