package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBooking(t *testing.T) {
	booking := Booking{
		ID:            "booking123",
		UserID:        "user123",
		HotelID:       "hotel123",
		RoomID:        "room123",
		CheckInDate:   time.Now(),
		CheckOutDate:  time.Now().AddDate(0, 0, 2),
		TotalPrice:    10000.0,
		Status:        "confirmed",
		PaymentStatus: "paid",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	assert.Equal(t, "booking123", booking.ID)
	assert.Equal(t, "user123", booking.UserID)
	assert.Equal(t, "hotel123", booking.HotelID)
	assert.Equal(t, "room123", booking.RoomID)
	assert.Equal(t, 10000.0, booking.TotalPrice)
	assert.Equal(t, "confirmed", booking.Status)
	assert.Equal(t, "paid", booking.PaymentStatus)
}

func TestBookingEvent(t *testing.T) {
	event := BookingEvent{
		BookingID:    "booking123",
		UserID:       "user123",
		HotelID:      "hotel123",
		RoomID:       "room123",
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().AddDate(0, 0, 2),
		TotalPrice:   10000.0,
		EventType:    "booking.created",
		Timestamp:    time.Now(),
	}

	assert.Equal(t, "booking123", event.BookingID)
	assert.Equal(t, "user123", event.UserID)
	assert.Equal(t, "hotel123", event.HotelID)
	assert.Equal(t, "booking.created", event.EventType)
	assert.Equal(t, 10000.0, event.TotalPrice)
}
