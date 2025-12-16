package service

import (
	"context"
	"fmt"
	"hotel-booking-system/internal/booking/domain"
	"hotel-booking-system/pkg/httpclient"
	"hotel-booking-system/pkg/logger"
)

type DeliveryClient interface {
	SendNotification(ctx context.Context, req *httpclient.SendNotificationRequest) error
}

type HotelClient interface {
	GetHotelOwnerID(ctx context.Context, hotelID string) (string, error)
}

type NotificationService struct {
	deliveryClient DeliveryClient
	hotelClient    HotelClient
}

func NewNotificationService(deliveryClient DeliveryClient, hotelClient HotelClient) *NotificationService {
	return &NotificationService{
		deliveryClient: deliveryClient,
		hotelClient:    hotelClient,
	}
}

func (ns *NotificationService) ProcessBookingEvent(ctx context.Context, event domain.BookingEvent) error {
	clientMessage := FormatBookingNotificationForClient(
		event.BookingID,
		event.HotelID,
		event.TotalPrice,
		event.CheckInDate,
		event.CheckOutDate,
	)

	if err := ns.deliveryClient.SendNotification(ctx, &httpclient.SendNotificationRequest{
		Channel:   "email",
		Recipient: event.UserID,
		Subject:   "Бронирование подтверждено",
		Message:   clientMessage,
	}); err != nil {
		logger.GetLogger().WithError(err).Error("failed to send notification to client")
	}

	ownerID, err := ns.hotelClient.GetHotelOwnerID(ctx, event.HotelID)
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to get hotel owner ID")
	} else {
		hotelierMessage := FormatBookingNotificationForHotelier(
			event.BookingID,
			event.UserID,
			event.HotelID,
			event.TotalPrice,
			event.CheckInDate,
			event.CheckOutDate,
		)

		if err := ns.deliveryClient.SendNotification(ctx, &httpclient.SendNotificationRequest{
			Channel:   "email",
			Recipient: ownerID,
			Subject:   "Новое бронирование в вашем отеле",
			Message:   hotelierMessage,
		}); err != nil {
			logger.GetLogger().WithError(err).Error("failed to send notification to hotelier")
		}
	}

	return nil
}

func FormatBookingNotificationForClient(bookingID, hotelID string, totalPrice float64, checkIn, checkOut interface{}) string {
	return fmt.Sprintf(
		"Ваше бронирование подтверждено!\n\nID бронирования: %s\nОтель: %s\nСумма: %.2f руб.\nДата заезда: %v\nДата выезда: %v\n\nСпасибо за выбор нашего сервиса!",
		bookingID, hotelID, totalPrice, checkIn, checkOut,
	)
}

func FormatBookingNotificationForHotelier(bookingID, userID, hotelID string, totalPrice float64, checkIn, checkOut interface{}) string {
	return fmt.Sprintf(
		"Новое бронирование в вашем отеле!\n\nID бронирования: %s\nПользователь: %s\nОтель: %s\nСумма: %.2f руб.\nДата заезда: %v\nДата выезда: %v",
		bookingID, userID, hotelID, totalPrice, checkIn, checkOut,
	)
}
