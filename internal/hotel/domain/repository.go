package domain

import "context"

type HotelRepository interface {
	CreateHotel(ctx context.Context, hotel *Hotel) error
	GetHotelByID(ctx context.Context, id string) (*Hotel, error)
	GetHotels(ctx context.Context, limit, offset int) ([]Hotel, error)
	GetHotelsByOwner(ctx context.Context, ownerID string) ([]Hotel, error)
	UpdateHotel(ctx context.Context, hotel *Hotel) error
	DeleteHotel(ctx context.Context, id string) error
}

type RoomRepository interface {
	CreateRoom(ctx context.Context, room *Room) error
	GetRoomByID(ctx context.Context, id string) (*Room, error)
	GetRoomsByHotel(ctx context.Context, hotelID string) ([]Room, error)
	UpdateRoom(ctx context.Context, room *Room) error
	DeleteRoom(ctx context.Context, id string) error
	GetRoomPrice(ctx context.Context, hotelID, roomID string) (float64, error)
}

type HotelUseCase interface {
	CreateHotel(ctx context.Context, hotel *Hotel) error
	GetHotel(ctx context.Context, id string) (*Hotel, error)
	GetHotels(ctx context.Context, limit, offset int) ([]Hotel, error)
	UpdateHotel(ctx context.Context, hotel *Hotel) error
	CreateRoom(ctx context.Context, room *Room) error
	GetHotelWithRooms(ctx context.Context, hotelID string) (*HotelWithRooms, error)
}
