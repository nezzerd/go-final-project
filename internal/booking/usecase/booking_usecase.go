package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"hotel-booking-system/internal/booking/domain"

	"github.com/google/uuid"
)

type HotelClient interface {
	GetRoomPrice(ctx context.Context, hotelID, roomID string) (float64, error)
}

type MessageProducer interface {
	SendMessage(ctx context.Context, key string, value interface{}) error
}

type PaymentClient interface {
	CreatePayment(ctx context.Context, bookingID string, amount float64) error
}

type BookingUseCase struct {
	repo          domain.BookingRepository
	hotelClient   HotelClient
	producer      MessageProducer
	paymentClient PaymentClient
}

func NewBookingUseCase(repo domain.BookingRepository, hotelClient HotelClient, producer MessageProducer, paymentClient PaymentClient) *BookingUseCase {
	return &BookingUseCase{
		repo:          repo,
		hotelClient:   hotelClient,
		producer:      producer,
		paymentClient: paymentClient,
	}
}

func (uc *BookingUseCase) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	if booking.CheckInDate.After(booking.CheckOutDate) {
		return errors.New("check-in date must be before check-out date")
	}

	pricePerNight, err := uc.hotelClient.GetRoomPrice(ctx, booking.HotelID, booking.RoomID)
	if err != nil {
		return err
	}

	nights := int(booking.CheckOutDate.Sub(booking.CheckInDate).Hours() / 24)
	if nights < 1 {
		nights = 1
	}
	booking.TotalPrice = pricePerNight * float64(nights)

	booking.ID = uuid.New().String()
	booking.Status = "confirmed"
	booking.PaymentStatus = "pending"

	if err := uc.repo.CreateBooking(ctx, booking); err != nil {
		return err
	}

	if uc.paymentClient != nil {
		if err := uc.paymentClient.CreatePayment(ctx, booking.ID, booking.TotalPrice); err != nil {
			return err
		}
	}

	event := domain.BookingEvent{
		BookingID:    booking.ID,
		UserID:       booking.UserID,
		HotelID:      booking.HotelID,
		RoomID:       booking.RoomID,
		CheckInDate:  booking.CheckInDate,
		CheckOutDate: booking.CheckOutDate,
		TotalPrice:   booking.TotalPrice,
		EventType:    "booking.created",
		Timestamp:    time.Now(),
	}

	if err := uc.producer.SendMessage(ctx, booking.ID, event); err != nil {
		return err
	}

	return nil
}

func (uc *BookingUseCase) GetBooking(ctx context.Context, id string) (*domain.Booking, error) {
	return uc.repo.GetBookingByID(ctx, id)
}

func (uc *BookingUseCase) GetBookingsByUser(ctx context.Context, userID string) ([]domain.Booking, error) {
	return uc.repo.GetBookingsByUser(ctx, userID)
}

func (uc *BookingUseCase) GetBookingsByHotel(ctx context.Context, hotelID string) ([]domain.Booking, error) {
	return uc.repo.GetBookingsByHotel(ctx, hotelID)
}

func (uc *BookingUseCase) UpdatePaymentStatus(ctx context.Context, id, status string) error {
	validStatuses := []string{"pending", "paid", "failed", "refunded"}
	found := false
	for _, s := range validStatuses {
		if strings.EqualFold(s, status) {
			found = true
			break
		}
	}
	if !found {
		return errors.New("invalid payment status")
	}
	return uc.repo.UpdatePaymentStatus(ctx, id, status)
}
