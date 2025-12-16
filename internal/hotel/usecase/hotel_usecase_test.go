package usecase

import (
	"context"
	"errors"
	"testing"

	"hotel-booking-system/internal/hotel/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHotelRepository struct {
	mock.Mock
}

func (m *MockHotelRepository) CreateHotel(ctx context.Context, hotel *domain.Hotel) error {
	args := m.Called(ctx, hotel)
	return args.Error(0)
}

func (m *MockHotelRepository) GetHotelByID(ctx context.Context, id string) (*domain.Hotel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Hotel), args.Error(1)
}

func (m *MockHotelRepository) GetHotels(ctx context.Context, limit, offset int) ([]domain.Hotel, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Hotel), args.Error(1)
}

func (m *MockHotelRepository) GetHotelsByOwner(ctx context.Context, ownerID string) ([]domain.Hotel, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).([]domain.Hotel), args.Error(1)
}

func (m *MockHotelRepository) UpdateHotel(ctx context.Context, hotel *domain.Hotel) error {
	args := m.Called(ctx, hotel)
	return args.Error(0)
}

func (m *MockHotelRepository) DeleteHotel(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) CreateRoom(ctx context.Context, room *domain.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockRoomRepository) GetRoomByID(ctx context.Context, id string) (*domain.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Room), args.Error(1)
}

func (m *MockRoomRepository) GetRoomsByHotel(ctx context.Context, hotelID string) ([]domain.Room, error) {
	args := m.Called(ctx, hotelID)
	return args.Get(0).([]domain.Room), args.Error(1)
}

func (m *MockRoomRepository) UpdateRoom(ctx context.Context, room *domain.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockRoomRepository) DeleteRoom(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoomRepository) GetRoomPrice(ctx context.Context, hotelID, roomID string) (float64, error) {
	args := m.Called(ctx, hotelID, roomID)
	return args.Get(0).(float64), args.Error(1)
}

func TestCreateHotel_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	hotel := &domain.Hotel{
		Name:    "Test Hotel",
		Address: "Test Address",
		OwnerID: "owner123",
	}

	mockHotelRepo.On("CreateHotel", mock.Anything, mock.Anything).Return(nil)

	err := uc.CreateHotel(context.Background(), hotel)
	assert.NoError(t, err)
	assert.NotEmpty(t, hotel.ID)
	mockHotelRepo.AssertExpectations(t)
}

func TestCreateHotel_InvalidData(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	hotel := &domain.Hotel{
		Name: "",
	}

	err := uc.CreateHotel(context.Background(), hotel)
	assert.Error(t, err)
}

func TestGetHotel_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	expectedHotel := &domain.Hotel{
		ID:      "hotel123",
		Name:    "Test Hotel",
		Address: "Test Address",
	}

	mockHotelRepo.On("GetHotelByID", mock.Anything, "hotel123").Return(expectedHotel, nil)

	hotel, err := uc.GetHotel(context.Background(), "hotel123")
	assert.NoError(t, err)
	assert.Equal(t, expectedHotel, hotel)
	mockHotelRepo.AssertExpectations(t)
}

func TestGetHotels_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	expectedHotels := []domain.Hotel{
		{ID: "hotel1", Name: "Hotel 1"},
		{ID: "hotel2", Name: "Hotel 2"},
	}

	mockHotelRepo.On("GetHotels", mock.Anything, 20, 0).Return(expectedHotels, nil)

	hotels, err := uc.GetHotels(context.Background(), 0, 0)
	assert.NoError(t, err)
	assert.Len(t, hotels, 2)
	mockHotelRepo.AssertExpectations(t)
}

func TestUpdateHotel_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	existingHotel := &domain.Hotel{
		ID:      "hotel123",
		OwnerID: "owner123",
	}

	updateHotel := &domain.Hotel{
		ID:      "hotel123",
		OwnerID: "owner123",
		Name:    "Updated Hotel",
	}

	mockHotelRepo.On("GetHotelByID", mock.Anything, "hotel123").Return(existingHotel, nil)
	mockHotelRepo.On("UpdateHotel", mock.Anything, updateHotel).Return(nil)

	err := uc.UpdateHotel(context.Background(), updateHotel)
	assert.NoError(t, err)
	mockHotelRepo.AssertExpectations(t)
}

func TestUpdateHotel_Unauthorized(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	existingHotel := &domain.Hotel{
		ID:      "hotel123",
		OwnerID: "owner123",
	}

	updateHotel := &domain.Hotel{
		ID:      "hotel123",
		OwnerID: "different_owner",
	}

	mockHotelRepo.On("GetHotelByID", mock.Anything, "hotel123").Return(existingHotel, nil)

	err := uc.UpdateHotel(context.Background(), updateHotel)
	assert.Error(t, err)
	mockHotelRepo.AssertExpectations(t)
}

func TestCreateRoom_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	room := &domain.Room{
		HotelID:       "hotel123",
		RoomNumber:    "101",
		PricePerNight: 5000,
	}

	mockRoomRepo.On("CreateRoom", mock.Anything, mock.Anything).Return(nil)

	err := uc.CreateRoom(context.Background(), room)
	assert.NoError(t, err)
	assert.NotEmpty(t, room.ID)
	mockRoomRepo.AssertExpectations(t)
}

func TestGetRoomPrice_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	mockRoomRepo.On("GetRoomPrice", mock.Anything, "hotel123", "room123").Return(5000.0, nil)

	price, err := uc.GetRoomPrice(context.Background(), "hotel123", "room123")
	assert.NoError(t, err)
	assert.Equal(t, 5000.0, price)
	mockRoomRepo.AssertExpectations(t)
}

func TestGetHotelWithRooms_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	hotel := &domain.Hotel{
		ID:   "hotel123",
		Name: "Test Hotel",
	}

	rooms := []domain.Room{
		{ID: "room1", HotelID: "hotel123"},
		{ID: "room2", HotelID: "hotel123"},
	}

	mockHotelRepo.On("GetHotelByID", mock.Anything, "hotel123").Return(hotel, nil)
	mockRoomRepo.On("GetRoomsByHotel", mock.Anything, "hotel123").Return(rooms, nil)

	result, err := uc.GetHotelWithRooms(context.Background(), "hotel123")
	assert.NoError(t, err)
	assert.Equal(t, hotel.ID, result.Hotel.ID)
	assert.Len(t, result.Rooms, 2)
	mockHotelRepo.AssertExpectations(t)
	mockRoomRepo.AssertExpectations(t)
}

func TestGetHotelWithRooms_HotelNotFound(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	mockHotelRepo.On("GetHotelByID", mock.Anything, "hotel123").Return(nil, errors.New("not found"))

	result, err := uc.GetHotelWithRooms(context.Background(), "hotel123")
	assert.Error(t, err)
	assert.Nil(t, result)
	mockHotelRepo.AssertExpectations(t)
}

func TestGetHotelsByOwner_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	expectedHotels := []domain.Hotel{
		{ID: "hotel1", OwnerID: "owner123"},
		{ID: "hotel2", OwnerID: "owner123"},
	}

	mockHotelRepo.On("GetHotelsByOwner", mock.Anything, "owner123").Return(expectedHotels, nil)

	hotels, err := uc.GetHotelsByOwner(context.Background(), "owner123")
	assert.NoError(t, err)
	assert.Len(t, hotels, 2)
	mockHotelRepo.AssertExpectations(t)
}

func TestDeleteHotel_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	hotel := &domain.Hotel{
		ID:      "hotel123",
		OwnerID: "owner123",
	}

	mockHotelRepo.On("GetHotelByID", mock.Anything, "hotel123").Return(hotel, nil)
	mockHotelRepo.On("DeleteHotel", mock.Anything, "hotel123").Return(nil)

	err := uc.DeleteHotel(context.Background(), "hotel123", "owner123")
	assert.NoError(t, err)
	mockHotelRepo.AssertExpectations(t)
}

func TestGetRoom_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	expectedRoom := &domain.Room{
		ID:       "room123",
		HotelID:  "hotel123",
		RoomType: "Standard",
	}

	mockRoomRepo.On("GetRoomByID", mock.Anything, "room123").Return(expectedRoom, nil)

	room, err := uc.GetRoom(context.Background(), "room123")
	assert.NoError(t, err)
	assert.Equal(t, expectedRoom, room)
	mockRoomRepo.AssertExpectations(t)
}

func TestGetRoomsByHotel_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	expectedRooms := []domain.Room{
		{ID: "room1", HotelID: "hotel123"},
		{ID: "room2", HotelID: "hotel123"},
	}

	mockRoomRepo.On("GetRoomsByHotel", mock.Anything, "hotel123").Return(expectedRooms, nil)

	rooms, err := uc.GetRoomsByHotel(context.Background(), "hotel123")
	assert.NoError(t, err)
	assert.Len(t, rooms, 2)
	mockRoomRepo.AssertExpectations(t)
}

func TestUpdateRoom_Success(t *testing.T) {
	mockHotelRepo := new(MockHotelRepository)
	mockRoomRepo := new(MockRoomRepository)
	uc := NewHotelUseCase(mockHotelRepo, mockRoomRepo)

	room := &domain.Room{
		ID:            "room123",
		HotelID:       "hotel123",
		RoomNumber:    "101",
		PricePerNight: 6000,
	}

	mockRoomRepo.On("UpdateRoom", mock.Anything, room).Return(nil)

	err := uc.UpdateRoom(context.Background(), room)
	assert.NoError(t, err)
	mockRoomRepo.AssertExpectations(t)
}
