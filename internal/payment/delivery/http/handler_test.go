package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hotel-booking-system/internal/payment/domain"
	"hotel-booking-system/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) ProcessPayment(ctx context.Context, req *domain.PaymentRequest) (*domain.PaymentResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PaymentResponse), args.Error(1)
}

func TestPaymentHandler_CreatePayment(t *testing.T) {
	logger.Init("info")

	t.Run("success", func(t *testing.T) {
		mockService := new(MockPaymentService)
		mockService.On("ProcessPayment", mock.Anything, mock.AnythingOfType("*domain.PaymentRequest")).Return(
			&domain.PaymentResponse{
				PaymentID: "payment-123",
				Status:    "processing",
				Message:   "payment is being processed",
			},
			nil,
		)

		handler := NewPaymentHandler(mockService)

		reqBody := domain.PaymentRequest{
			BookingID: "booking-123",
			Amount:    1000.0,
			Currency:  "RUB",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/payments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreatePayment(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)
		mockService.AssertExpectations(t)

		var response domain.PaymentResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "payment-123", response.PaymentID)
		assert.Equal(t, "processing", response.Status)
	})

	t.Run("success without currency", func(t *testing.T) {
		mockService := new(MockPaymentService)
		mockService.On("ProcessPayment", mock.Anything, mock.MatchedBy(func(req *domain.PaymentRequest) bool {
			return req.Currency == "RUB"
		})).Return(
			&domain.PaymentResponse{
				PaymentID: "payment-456",
				Status:    "processing",
			},
			nil,
		)

		handler := NewPaymentHandler(mockService)

		reqBody := domain.PaymentRequest{
			BookingID: "booking-456",
			Amount:    2000.0,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/payments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreatePayment(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mockService := new(MockPaymentService)
		handler := NewPaymentHandler(mockService)

		req := httptest.NewRequest("POST", "/api/payments", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreatePayment(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "ProcessPayment")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockPaymentService)
		mockService.On("ProcessPayment", mock.Anything, mock.AnythingOfType("*domain.PaymentRequest")).Return(
			nil,
			assert.AnError,
		)

		handler := NewPaymentHandler(mockService)

		reqBody := domain.PaymentRequest{
			BookingID: "booking-789",
			Amount:    3000.0,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/payments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreatePayment(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestNewPaymentHandler(t *testing.T) {
	mockService := new(MockPaymentService)
	handler := NewPaymentHandler(mockService)
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.paymentService)
}
