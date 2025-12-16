package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"hotel-booking-system/internal/booking/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *MockBookingRepository) GetBookingByID(ctx context.Context, id string) (*domain.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) GetBookingsByUser(ctx context.Context, userID string) ([]domain.Booking, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) GetBookingsByHotel(ctx context.Context, hotelID string) ([]domain.Booking, error) {
	args := m.Called(ctx, hotelID)
	return args.Get(0).([]domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) UpdateBookingStatus(ctx context.Context, id, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockBookingRepository) UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) error {
	args := m.Called(ctx, id, paymentStatus)
	return args.Error(0)
}

type MockHotelClient struct {
	GetRoomPriceFunc func(ctx context.Context, hotelID, roomID string) (float64, error)
}

func (m *MockHotelClient) GetRoomPrice(ctx context.Context, hotelID, roomID string) (float64, error) {
	if m.GetRoomPriceFunc != nil {
		return m.GetRoomPriceFunc(ctx, hotelID, roomID)
	}
	return 0, nil
}

func (m *MockHotelClient) Close() error {
	return nil
}

type MockProducer struct {
	SendMessageFunc func(ctx context.Context, key string, value interface{}) error
}

func (m *MockProducer) SendMessage(ctx context.Context, key string, value interface{}) error {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(ctx, key, value)
	}
	return nil
}

func (m *MockProducer) Close() error {
	return nil
}

func TestCreateBooking_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	mockClient := &MockHotelClient{
		GetRoomPriceFunc: func(ctx context.Context, hotelID, roomID string) (float64, error) {
			return 5000.0, nil
		},
	}
	mockProducer := &MockProducer{
		SendMessageFunc: func(ctx context.Context, key string, value interface{}) error {
			return nil
		},
	}

	booking := &domain.Booking{
		UserID:       "user123",
		HotelID:      "hotel123",
		RoomID:       "room123",
		CheckInDate:  time.Now().AddDate(0, 0, 1),
		CheckOutDate: time.Now().AddDate(0, 0, 3),
	}

	mockRepo.On("CreateBooking", mock.Anything, mock.Anything).Return(nil)

	uc := &BookingUseCase{
		repo:        mockRepo,
		hotelClient: mockClient,
		producer:    mockProducer,
	}

	err := uc.CreateBooking(context.Background(), booking)
	assert.NoError(t, err)
	assert.NotEmpty(t, booking.ID)
	assert.Equal(t, "confirmed", booking.Status)
	assert.Equal(t, "pending", booking.PaymentStatus)
	assert.Greater(t, booking.TotalPrice, 0.0)
	mockRepo.AssertExpectations(t)
}

func TestCreateBooking_InvalidDates(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	mockClient := &MockHotelClient{}
	mockProducer := &MockProducer{}

	booking := &domain.Booking{
		UserID:       "user123",
		HotelID:      "hotel123",
		RoomID:       "room123",
		CheckInDate:  time.Now().AddDate(0, 0, 3),
		CheckOutDate: time.Now().AddDate(0, 0, 1),
	}

	uc := &BookingUseCase{
		repo:        mockRepo,
		hotelClient: mockClient,
		producer:    mockProducer,
	}

	err := uc.CreateBooking(context.Background(), booking)
	assert.Error(t, err)
}

func TestGetBooking_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	mockClient := &MockHotelClient{}
	mockProducer := &MockProducer{}

	expectedBooking := &domain.Booking{
		ID:      "booking123",
		UserID:  "user123",
		HotelID: "hotel123",
	}

	mockRepo.On("GetBookingByID", mock.Anything, "booking123").Return(expectedBooking, nil)

	uc := &BookingUseCase{
		repo:        mockRepo,
		hotelClient: mockClient,
		producer:    mockProducer,
	}

	booking, err := uc.GetBooking(context.Background(), "booking123")
	assert.NoError(t, err)
	assert.Equal(t, expectedBooking, booking)
	mockRepo.AssertExpectations(t)
}

func TestGetBookingsByUser_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	mockClient := &MockHotelClient{}
	mockProducer := &MockProducer{}

	expectedBookings := []domain.Booking{
		{ID: "booking1", UserID: "user123"},
		{ID: "booking2", UserID: "user123"},
	}

	mockRepo.On("GetBookingsByUser", mock.Anything, "user123").Return(expectedBookings, nil)

	uc := &BookingUseCase{
		repo:        mockRepo,
		hotelClient: mockClient,
		producer:    mockProducer,
	}

	bookings, err := uc.GetBookingsByUser(context.Background(), "user123")
	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetBookingsByHotel_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	mockClient := &MockHotelClient{}
	mockProducer := &MockProducer{}

	expectedBookings := []domain.Booking{
		{ID: "booking1", HotelID: "hotel123"},
		{ID: "booking2", HotelID: "hotel123"},
	}

	mockRepo.On("GetBookingsByHotel", mock.Anything, "hotel123").Return(expectedBookings, nil)

	uc := &BookingUseCase{
		repo:        mockRepo,
		hotelClient: mockClient,
		producer:    mockProducer,
	}

	bookings, err := uc.GetBookingsByHotel(context.Background(), "hotel123")
	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePaymentStatus_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	mockClient := &MockHotelClient{}
	mockProducer := &MockProducer{}

	mockRepo.On("UpdatePaymentStatus", mock.Anything, "booking123", "paid").Return(nil)

	uc := &BookingUseCase{
		repo:        mockRepo,
		hotelClient: mockClient,
		producer:    mockProducer,
	}

	err := uc.UpdatePaymentStatus(context.Background(), "booking123", "paid")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePaymentStatus_InvalidStatus(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	mockClient := &MockHotelClient{}
	mockProducer := &MockProducer{}

	uc := &BookingUseCase{
		repo:        mockRepo,
		hotelClient: mockClient,
		producer:    mockProducer,
	}

	err := uc.UpdatePaymentStatus(context.Background(), "booking123", "invalid")
	assert.Error(t, err)
}

func TestGetBooking_NotFound(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	mockClient := &MockHotelClient{}
	mockProducer := &MockProducer{}

	mockRepo.On("GetBookingByID", mock.Anything, "booking123").Return(nil, errors.New("not found"))

	uc := &BookingUseCase{
		repo:        mockRepo,
		hotelClient: mockClient,
		producer:    mockProducer,
	}

	booking, err := uc.GetBooking(context.Background(), "booking123")
	assert.Error(t, err)
	assert.Nil(t, booking)
	mockRepo.AssertExpectations(t)
}
