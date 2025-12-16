package usecase

import (
	"context"
	"errors"

	"hotel-booking-system/internal/hotel/domain"

	"github.com/google/uuid"
)

type HotelUseCase struct {
	hotelRepo domain.HotelRepository
	roomRepo  domain.RoomRepository
}

func NewHotelUseCase(hotelRepo domain.HotelRepository, roomRepo domain.RoomRepository) *HotelUseCase {
	return &HotelUseCase{
		hotelRepo: hotelRepo,
		roomRepo:  roomRepo,
	}
}

func (uc *HotelUseCase) CreateHotel(ctx context.Context, hotel *domain.Hotel) error {
	if hotel.Name == "" || hotel.Address == "" || hotel.OwnerID == "" {
		return errors.New("invalid hotel data")
	}
	hotel.ID = uuid.New().String()
	return uc.hotelRepo.CreateHotel(ctx, hotel)
}

func (uc *HotelUseCase) GetHotel(ctx context.Context, id string) (*domain.Hotel, error) {
	return uc.hotelRepo.GetHotelByID(ctx, id)
}

func (uc *HotelUseCase) GetHotels(ctx context.Context, limit, offset int) ([]domain.Hotel, error) {
	if limit <= 0 {
		limit = 20
	}
	return uc.hotelRepo.GetHotels(ctx, limit, offset)
}

func (uc *HotelUseCase) GetHotelsByOwner(ctx context.Context, ownerID string) ([]domain.Hotel, error) {
	return uc.hotelRepo.GetHotelsByOwner(ctx, ownerID)
}

func (uc *HotelUseCase) UpdateHotel(ctx context.Context, hotel *domain.Hotel) error {
	existing, err := uc.hotelRepo.GetHotelByID(ctx, hotel.ID)
	if err != nil {
		return err
	}
	if existing.OwnerID != hotel.OwnerID {
		return errors.New("unauthorized to update this hotel")
	}
	return uc.hotelRepo.UpdateHotel(ctx, hotel)
}

func (uc *HotelUseCase) DeleteHotel(ctx context.Context, id, ownerID string) error {
	existing, err := uc.hotelRepo.GetHotelByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.OwnerID != ownerID {
		return errors.New("unauthorized to delete this hotel")
	}
	return uc.hotelRepo.DeleteHotel(ctx, id)
}

func (uc *HotelUseCase) CreateRoom(ctx context.Context, room *domain.Room) error {
	room.ID = uuid.New().String()
	return uc.roomRepo.CreateRoom(ctx, room)
}

func (uc *HotelUseCase) GetRoom(ctx context.Context, id string) (*domain.Room, error) {
	return uc.roomRepo.GetRoomByID(ctx, id)
}

func (uc *HotelUseCase) GetRoomsByHotel(ctx context.Context, hotelID string) ([]domain.Room, error) {
	return uc.roomRepo.GetRoomsByHotel(ctx, hotelID)
}

func (uc *HotelUseCase) GetHotelWithRooms(ctx context.Context, hotelID string) (*domain.HotelWithRooms, error) {
	hotel, err := uc.hotelRepo.GetHotelByID(ctx, hotelID)
	if err != nil {
		return nil, err
	}

	rooms, err := uc.roomRepo.GetRoomsByHotel(ctx, hotelID)
	if err != nil {
		return nil, err
	}

	return &domain.HotelWithRooms{
		Hotel: *hotel,
		Rooms: rooms,
	}, nil
}

func (uc *HotelUseCase) UpdateRoom(ctx context.Context, room *domain.Room) error {
	return uc.roomRepo.UpdateRoom(ctx, room)
}

func (uc *HotelUseCase) GetRoomPrice(ctx context.Context, hotelID, roomID string) (float64, error) {
	return uc.roomRepo.GetRoomPrice(ctx, hotelID, roomID)
}
