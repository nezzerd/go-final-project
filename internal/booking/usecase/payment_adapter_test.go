package usecase

import (
	"context"
	"errors"
	"testing"

	"hotel-booking-system/pkg/httpclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentClient struct {
	mock.Mock
}

func (m *MockPaymentClient) CreatePayment(ctx context.Context, req *httpclient.PaymentRequest) (*httpclient.PaymentResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*httpclient.PaymentResponse), args.Error(1)
}

func TestNewPaymentClientAdapter(t *testing.T) {
	mockClient := new(MockPaymentClient)
	adapter := NewPaymentClientAdapter(mockClient)

	assert.NotNil(t, adapter)
	assert.IsType(t, &paymentClientAdapter{}, adapter)
}

func TestPaymentClientAdapter_CreatePayment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockClient := new(MockPaymentClient)
		mockClient.On("CreatePayment", mock.Anything, mock.MatchedBy(func(req *httpclient.PaymentRequest) bool {
			return req.BookingID == "booking-123" &&
				req.Amount == 1000.0 &&
				req.Currency == "RUB"
		})).Return(
			&httpclient.PaymentResponse{
				PaymentID: "payment-123",
				Status:    "processing",
			},
			nil,
		)

		adapter := NewPaymentClientAdapter(mockClient)

		err := adapter.CreatePayment(context.Background(), "booking-123", 1000.0)
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("client error", func(t *testing.T) {
		mockClient := new(MockPaymentClient)
		mockClient.On("CreatePayment", mock.Anything, mock.Anything).Return(
			nil,
			errors.New("payment service error"),
		)

		adapter := NewPaymentClientAdapter(mockClient)

		err := adapter.CreatePayment(context.Background(), "booking-123", 1000.0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment service error")
		mockClient.AssertExpectations(t)
	})

	t.Run("zero amount", func(t *testing.T) {
		mockClient := new(MockPaymentClient)
		mockClient.On("CreatePayment", mock.Anything, mock.MatchedBy(func(req *httpclient.PaymentRequest) bool {
			return req.Amount == 0.0
		})).Return(
			&httpclient.PaymentResponse{
				PaymentID: "payment-456",
				Status:    "processing",
			},
			nil,
		)

		adapter := NewPaymentClientAdapter(mockClient)

		err := adapter.CreatePayment(context.Background(), "booking-456", 0.0)
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("large amount", func(t *testing.T) {
		mockClient := new(MockPaymentClient)
		mockClient.On("CreatePayment", mock.Anything, mock.MatchedBy(func(req *httpclient.PaymentRequest) bool {
			return req.Amount == 999999.99
		})).Return(
			&httpclient.PaymentResponse{
				PaymentID: "payment-789",
				Status:    "processing",
			},
			nil,
		)

		adapter := NewPaymentClientAdapter(mockClient)

		err := adapter.CreatePayment(context.Background(), "booking-789", 999999.99)
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})
}
