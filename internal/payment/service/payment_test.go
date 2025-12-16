package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hotel-booking-system/internal/payment/domain"
	"hotel-booking-system/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPaymentService(t *testing.T) {
	logger.Init("info")

	webhookURL := "http://example.com/webhook"
	service := NewPaymentService(webhookURL)

	assert.NotNil(t, service)
	assert.Equal(t, webhookURL, service.webhookURL)
}

func TestPaymentService_ProcessPayment(t *testing.T) {
	logger.Init("info")

	t.Run("success with positive amount", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		service := NewPaymentService(server.URL + "/webhook")

		req := &domain.PaymentRequest{
			BookingID: "booking-123",
			Amount:    1000.0,
			Currency:  "RUB",
		}

		response, err := service.ProcessPayment(context.Background(), req)

		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.PaymentID)
		assert.Equal(t, "processing", response.Status)

		time.Sleep(3 * time.Second)
	})

	t.Run("success with zero amount", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		service := NewPaymentService(server.URL + "/webhook")

		req := &domain.PaymentRequest{
			BookingID: "booking-456",
			Amount:    0.0,
			Currency:  "RUB",
		}

		response, err := service.ProcessPayment(context.Background(), req)

		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "processing", response.Status)

		time.Sleep(3 * time.Second)
	})

	t.Run("success with negative amount", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		service := NewPaymentService(server.URL + "/webhook")

		req := &domain.PaymentRequest{
			BookingID: "booking-789",
			Amount:    -100.0,
			Currency:  "RUB",
		}

		response, err := service.ProcessPayment(context.Background(), req)

		require.NoError(t, err)
		assert.NotNil(t, response)

		time.Sleep(3 * time.Second)
	})
}

func TestPaymentService_sendWebhook(t *testing.T) {
	logger.Init("info")

	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		service := NewPaymentService(server.URL + "/webhook")

		webhook := domain.PaymentWebhook{
			PaymentID:   "payment-123",
			BookingID:   "booking-123",
			Status:      "paid",
			Amount:      1000.0,
			ProcessedAt: time.Now().Format(time.RFC3339),
		}

		err := service.sendWebhook(context.Background(), webhook)
		assert.NoError(t, err)
	})

	t.Run("webhook returns error status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		service := NewPaymentService(server.URL + "/webhook")

		webhook := domain.PaymentWebhook{
			PaymentID:   "payment-123",
			BookingID:   "booking-123",
			Status:      "paid",
			Amount:      1000.0,
			ProcessedAt: time.Now().Format(time.RFC3339),
		}

		err := service.sendWebhook(context.Background(), webhook)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status 500")
	})

	t.Run("invalid webhook URL", func(t *testing.T) {
		service := NewPaymentService("http://invalid-url-that-does-not-exist:9999/webhook")

		webhook := domain.PaymentWebhook{
			PaymentID:   "payment-123",
			BookingID:   "booking-123",
			Status:      "paid",
			Amount:      1000.0,
			ProcessedAt: time.Now().Format(time.RFC3339),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := service.sendWebhook(ctx, webhook)
		assert.Error(t, err)
	})
}
