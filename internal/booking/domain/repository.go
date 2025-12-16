package domain

import "context"

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *Booking) error
	GetBookingByID(ctx context.Context, id string) (*Booking, error)
	GetBookingsByUser(ctx context.Context, userID string) ([]Booking, error)
	GetBookingsByHotel(ctx context.Context, hotelID string) ([]Booking, error)
	UpdateBookingStatus(ctx context.Context, id, status string) error
	UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) error
}

type BookingUseCase interface {
	CreateBooking(ctx context.Context, booking *Booking) error
	GetBooking(ctx context.Context, id string) (*Booking, error)
	GetBookingsByUser(ctx context.Context, userID string) ([]Booking, error)
	GetBookingsByHotel(ctx context.Context, hotelID string) ([]Booking, error)
	UpdatePaymentStatus(ctx context.Context, id, status string) error
}
