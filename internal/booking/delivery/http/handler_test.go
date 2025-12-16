package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hotel-booking-system/internal/booking/domain"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookingUseCase struct {
	mock.Mock
}

func (m *MockBookingUseCase) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *MockBookingUseCase) GetBooking(ctx context.Context, id string) (*domain.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingUseCase) GetBookingsByUser(ctx context.Context, userID string) ([]domain.Booking, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Booking), args.Error(1)
}

func (m *MockBookingUseCase) GetBookingsByHotel(ctx context.Context, hotelID string) ([]domain.Booking, error) {
	args := m.Called(ctx, hotelID)
	return args.Get(0).([]domain.Booking), args.Error(1)
}

func (m *MockBookingUseCase) UpdatePaymentStatus(ctx context.Context, id, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestCreateBooking_Success(t *testing.T) {
	mockUC := new(MockBookingUseCase)
	handler := NewBookingHandler(mockUC)

	booking := domain.Booking{
		UserID:  "user123",
		HotelID: "hotel123",
		RoomID:  "room123",
	}

	mockUC.On("CreateBooking", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(booking)
	req := httptest.NewRequest("POST", "/api/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetBooking_Success(t *testing.T) {
	mockUC := new(MockBookingUseCase)
	handler := NewBookingHandler(mockUC)

	booking := &domain.Booking{
		ID:      "booking123",
		UserID:  "user123",
		HotelID: "hotel123",
	}

	mockUC.On("GetBooking", mock.Anything, "booking123").Return(booking, nil)

	req := httptest.NewRequest("GET", "/api/bookings/booking123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "booking123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetBooking(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetBookingsByUser_Success(t *testing.T) {
	mockUC := new(MockBookingUseCase)
	handler := NewBookingHandler(mockUC)

	bookings := []domain.Booking{
		{ID: "1", UserID: "user123"},
		{ID: "2", UserID: "user123"},
	}

	mockUC.On("GetBookingsByUser", mock.Anything, "user123").Return(bookings, nil)

	req := httptest.NewRequest("GET", "/api/bookings/user/user123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userId", "user123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetBookingsByUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetBookingsByHotel_Success(t *testing.T) {
	mockUC := new(MockBookingUseCase)
	handler := NewBookingHandler(mockUC)

	bookings := []domain.Booking{
		{ID: "1", HotelID: "hotel123"},
		{ID: "2", HotelID: "hotel123"},
	}

	mockUC.On("GetBookingsByHotel", mock.Anything, "hotel123").Return(bookings, nil)

	req := httptest.NewRequest("GET", "/api/bookings/hotel/hotel123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("hotelId", "hotel123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetBookingsByHotel(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}
