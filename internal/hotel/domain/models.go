package domain

import (
	"time"
)

type Hotel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Room struct {
	ID            string    `json:"id"`
	HotelID       string    `json:"hotel_id"`
	RoomNumber    string    `json:"room_number"`
	RoomType      string    `json:"room_type"`
	PricePerNight float64   `json:"price_per_night"`
	Capacity      int       `json:"capacity"`
	Description   string    `json:"description"`
	IsAvailable   bool      `json:"is_available"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type HotelWithRooms struct {
	Hotel Hotel  `json:"hotel"`
	Rooms []Room `json:"rooms"`
}
