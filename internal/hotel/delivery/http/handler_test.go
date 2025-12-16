package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"hotel-booking-system/internal/hotel/domain"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHotelUseCase struct {
	mock.Mock
}

func (m *MockHotelUseCase) CreateHotel(ctx context.Context, hotel *domain.Hotel) error {
	args := m.Called(ctx, hotel)
	return args.Error(0)
}

func (m *MockHotelUseCase) GetHotel(ctx context.Context, id string) (*domain.Hotel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Hotel), args.Error(1)
}

func (m *MockHotelUseCase) GetHotels(ctx context.Context, limit, offset int) ([]domain.Hotel, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Hotel), args.Error(1)
}

func (m *MockHotelUseCase) UpdateHotel(ctx context.Context, hotel *domain.Hotel) error {
	args := m.Called(ctx, hotel)
	return args.Error(0)
}

func (m *MockHotelUseCase) CreateRoom(ctx context.Context, room *domain.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockHotelUseCase) GetHotelWithRooms(ctx context.Context, hotelID string) (*domain.HotelWithRooms, error) {
	args := m.Called(ctx, hotelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.HotelWithRooms), args.Error(1)
}

func TestCreateHotel_Success(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	hotel := domain.Hotel{
		Name:    "Test Hotel",
		Address: "Test Address",
		OwnerID: "owner123",
	}

	mockUC.On("CreateHotel", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(hotel)
	req := httptest.NewRequest("POST", "/api/hotels", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateHotel(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetHotels_Success(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	hotels := []domain.Hotel{
		{ID: "1", Name: "Hotel 1"},
		{ID: "2", Name: "Hotel 2"},
	}

	mockUC.On("GetHotels", mock.Anything, 20, 0).Return(hotels, nil)

	req := httptest.NewRequest("GET", "/api/hotels?limit=20&offset=0", nil)
	w := httptest.NewRecorder()

	handler.GetHotels(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetHotel_Success(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	hotel := &domain.Hotel{
		ID:   "hotel123",
		Name: "Test Hotel",
	}

	mockUC.On("GetHotel", mock.Anything, "hotel123").Return(hotel, nil)

	req := httptest.NewRequest("GET", "/api/hotels/hotel123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "hotel123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetHotel(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetHotelWithRooms_Success(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	hotelWithRooms := &domain.HotelWithRooms{
		Hotel: domain.Hotel{
			ID:   "hotel123",
			Name: "Test Hotel",
		},
		Rooms: []domain.Room{
			{ID: "room1", RoomNumber: "101"},
		},
	}

	mockUC.On("GetHotelWithRooms", mock.Anything, "hotel123").Return(hotelWithRooms, nil)

	req := httptest.NewRequest("GET", "/api/hotels/hotel123/rooms", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "hotel123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetHotelWithRooms(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateRoom_Success(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	room := domain.Room{
		HotelID:       "hotel123",
		RoomNumber:    "101",
		RoomType:      "Standard",
		PricePerNight: 5000,
	}

	mockUC.On("CreateRoom", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(room)
	req := httptest.NewRequest("POST", "/api/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateHotel_InvalidJSON(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	req := httptest.NewRequest("POST", "/api/hotels", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateHotel(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateHotel_UseCaseError(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	hotel := domain.Hotel{
		Name:    "Test Hotel",
		Address: "Test Address",
	}

	mockUC.On("CreateHotel", mock.Anything, mock.Anything).Return(errors.New("database error"))

	body, _ := json.Marshal(hotel)
	req := httptest.NewRequest("POST", "/api/hotels", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateHotel(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetHotels_Error(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	mockUC.On("GetHotels", mock.Anything, 20, 0).Return([]domain.Hotel{}, errors.New("database error"))

	req := httptest.NewRequest("GET", "/api/hotels?limit=20&offset=0", nil)
	w := httptest.NewRecorder()

	handler.GetHotels(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetHotelWithRooms_NotFound(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	mockUC.On("GetHotelWithRooms", mock.Anything, "hotel123").Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/api/hotels/hotel123/rooms", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "hotel123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetHotelWithRooms(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateRoom_InvalidJSON(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	req := httptest.NewRequest("POST", "/api/rooms", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateRoom_UseCaseError(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	room := domain.Room{
		HotelID:       "hotel123",
		RoomNumber:    "101",
		RoomType:      "Standard",
		PricePerNight: 5000,
	}

	mockUC.On("CreateRoom", mock.Anything, mock.Anything).Return(errors.New("database error"))

	body, _ := json.Marshal(room)
	req := httptest.NewRequest("POST", "/api/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestUpdateHotel_Success(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	hotel := domain.Hotel{
		ID:      "hotel123",
		Name:    "Updated Hotel",
		Address: "New Address",
	}

	mockUC.On("UpdateHotel", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(hotel)
	req := httptest.NewRequest("PUT", "/api/hotels/hotel123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "hotel123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UpdateHotel(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestUpdateHotel_InvalidJSON(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	req := httptest.NewRequest("PUT", "/api/hotels/hotel123", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "hotel123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UpdateHotel(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateHotel_UseCaseError(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	hotel := domain.Hotel{
		ID:      "hotel123",
		Name:    "Updated Hotel",
		Address: "New Address",
	}

	mockUC.On("UpdateHotel", mock.Anything, mock.Anything).Return(errors.New("update error"))

	body, _ := json.Marshal(hotel)
	req := httptest.NewRequest("PUT", "/api/hotels/hotel123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "hotel123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UpdateHotel(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}
