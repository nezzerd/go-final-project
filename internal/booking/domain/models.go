package domain

import (
	"time"
)

type Booking struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	HotelID       string    `json:"hotel_id"`
	RoomID        string    `json:"room_id"`
	CheckInDate   time.Time `json:"check_in_date"`
	CheckOutDate  time.Time `json:"check_out_date"`
	TotalPrice    float64   `json:"total_price"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BookingEvent struct {
	BookingID    string    `json:"booking_id"`
	UserID       string    `json:"user_id"`
	HotelID      string    `json:"hotel_id"`
	RoomID       string    `json:"room_id"`
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
	TotalPrice   float64   `json:"total_price"`
	EventType    string    `json:"event_type"`
	Timestamp    time.Time `json:"timestamp"`
}
