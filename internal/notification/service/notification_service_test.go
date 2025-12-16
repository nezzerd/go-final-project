package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"hotel-booking-system/internal/booking/domain"
	"hotel-booking-system/pkg/httpclient"
	"hotel-booking-system/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDeliveryClient struct {
	mock.Mock
}

func (m *MockDeliveryClient) SendNotification(ctx context.Context, req *httpclient.SendNotificationRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

type MockHotelClient struct {
	mock.Mock
}

func (m *MockHotelClient) GetHotelOwnerID(ctx context.Context, hotelID string) (string, error) {
	args := m.Called(ctx, hotelID)
	return args.String(0), args.Error(1)
}

func TestNewNotificationService(t *testing.T) {
	logger.Init("info")

	mockDeliveryClient := new(MockDeliveryClient)
	mockHotelClient := new(MockHotelClient)

	service := NewNotificationService(mockDeliveryClient, mockHotelClient)

	assert.NotNil(t, service)
	assert.Equal(t, mockDeliveryClient, service.deliveryClient)
	assert.Equal(t, mockHotelClient, service.hotelClient)
}

func TestNotificationService_ProcessBookingEvent(t *testing.T) {
	logger.Init("info")

	event := domain.BookingEvent{
		BookingID:    "booking-123",
		UserID:       "user-123",
		HotelID:      "hotel-123",
		RoomID:       "room-123",
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(24 * time.Hour),
		TotalPrice:   5000.0,
		EventType:    "booking.created",
		Timestamp:    time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockDeliveryClient := new(MockDeliveryClient)
		mockDeliveryClient.On("SendNotification", mock.Anything, mock.Anything).Return(nil).Twice()

		mockHotelClient := new(MockHotelClient)
		mockHotelClient.On("GetHotelOwnerID", mock.Anything, "hotel-123").Return("owner-123", nil)

		service := NewNotificationService(mockDeliveryClient, mockHotelClient)

		err := service.ProcessBookingEvent(context.Background(), event)

		assert.NoError(t, err)
		mockDeliveryClient.AssertExpectations(t)
		mockHotelClient.AssertExpectations(t)
	})

	t.Run("delivery client error for client", func(t *testing.T) {
		mockDeliveryClient := new(MockDeliveryClient)
		mockDeliveryClient.On("SendNotification", mock.Anything, mock.Anything).Return(errors.New("delivery error")).Once()

		mockHotelClient := new(MockHotelClient)
		mockHotelClient.On("GetHotelOwnerID", mock.Anything, "hotel-123").Return("owner-123", nil)
		mockDeliveryClient.On("SendNotification", mock.Anything, mock.Anything).Return(nil).Once()

		service := NewNotificationService(mockDeliveryClient, mockHotelClient)

		err := service.ProcessBookingEvent(context.Background(), event)

		assert.NoError(t, err)
		mockDeliveryClient.AssertNumberOfCalls(t, "SendNotification", 2)
	})

	t.Run("hotel client error", func(t *testing.T) {
		mockDeliveryClient := new(MockDeliveryClient)
		mockDeliveryClient.On("SendNotification", mock.Anything, mock.Anything).Return(nil).Once()

		mockHotelClient := new(MockHotelClient)
		mockHotelClient.On("GetHotelOwnerID", mock.Anything, "hotel-123").Return("", errors.New("hotel not found"))

		service := NewNotificationService(mockDeliveryClient, mockHotelClient)

		err := service.ProcessBookingEvent(context.Background(), event)

		assert.NoError(t, err)
		mockDeliveryClient.AssertNumberOfCalls(t, "SendNotification", 1)
		mockHotelClient.AssertExpectations(t)
	})

	t.Run("delivery client error for hotelier", func(t *testing.T) {
		mockDeliveryClient := new(MockDeliveryClient)
		mockDeliveryClient.On("SendNotification", mock.Anything, mock.Anything).Return(nil).Once()
		mockDeliveryClient.On("SendNotification", mock.Anything, mock.Anything).Return(errors.New("delivery error")).Once()

		mockHotelClient := new(MockHotelClient)
		mockHotelClient.On("GetHotelOwnerID", mock.Anything, "hotel-123").Return("owner-123", nil)

		service := NewNotificationService(mockDeliveryClient, mockHotelClient)

		err := service.ProcessBookingEvent(context.Background(), event)

		assert.NoError(t, err)
		mockDeliveryClient.AssertNumberOfCalls(t, "SendNotification", 2)
	})
}

func TestFormatBookingNotificationForClient(t *testing.T) {
	message := FormatBookingNotificationForClient(
		"booking-123",
		"hotel-123",
		5000.0,
		time.Date(2024, 12, 20, 14, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 25, 12, 0, 0, 0, time.UTC),
	)

	assert.Contains(t, message, "booking-123")
	assert.Contains(t, message, "hotel-123")
	assert.Contains(t, message, "5000.00")
}

func TestFormatBookingNotificationForHotelier(t *testing.T) {
	message := FormatBookingNotificationForHotelier(
		"booking-123",
		"user-123",
		"hotel-123",
		5000.0,
		time.Date(2024, 12, 20, 14, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 25, 12, 0, 0, 0, time.UTC),
	)

	assert.Contains(t, message, "booking-123")
	assert.Contains(t, message, "user-123")
	assert.Contains(t, message, "hotel-123")
	assert.Contains(t, message, "5000.00")
}
