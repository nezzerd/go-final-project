package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHotel(t *testing.T) {
	hotel := Hotel{
		ID:          "hotel123",
		Name:        "Test Hotel",
		Description: "Test Description",
		Address:     "Test Address",
		OwnerID:     "owner123",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	assert.Equal(t, "hotel123", hotel.ID)
	assert.Equal(t, "Test Hotel", hotel.Name)
	assert.Equal(t, "Test Description", hotel.Description)
	assert.Equal(t, "Test Address", hotel.Address)
	assert.Equal(t, "owner123", hotel.OwnerID)
}

func TestRoom(t *testing.T) {
	room := Room{
		ID:            "room123",
		HotelID:       "hotel123",
		RoomNumber:    "101",
		RoomType:      "Standard",
		PricePerNight: 5000.0,
		Capacity:      2,
		Description:   "Test Room",
		IsAvailable:   true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	assert.Equal(t, "room123", room.ID)
	assert.Equal(t, "hotel123", room.HotelID)
	assert.Equal(t, "101", room.RoomNumber)
	assert.Equal(t, "Standard", room.RoomType)
	assert.Equal(t, 5000.0, room.PricePerNight)
	assert.Equal(t, 2, room.Capacity)
	assert.True(t, room.IsAvailable)
}

func TestHotelWithRooms(t *testing.T) {
	hotel := Hotel{
		ID:   "hotel123",
		Name: "Test Hotel",
	}

	rooms := []Room{
		{ID: "room1", HotelID: "hotel123"},
		{ID: "room2", HotelID: "hotel123"},
	}

	hotelWithRooms := HotelWithRooms{
		Hotel: hotel,
		Rooms: rooms,
	}

	assert.Equal(t, hotel.ID, hotelWithRooms.Hotel.ID)
	assert.Len(t, hotelWithRooms.Rooms, 2)
}
