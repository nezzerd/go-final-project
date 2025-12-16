package usecase

import (
	"context"
	"hotel-booking-system/pkg/httpclient"
)

type PaymentClientInterface interface {
	CreatePayment(ctx context.Context, req *httpclient.PaymentRequest) (*httpclient.PaymentResponse, error)
}

type paymentClientAdapter struct {
	client PaymentClientInterface
}

func NewPaymentClientAdapter(client PaymentClientInterface) PaymentClient {
	return &paymentClientAdapter{client: client}
}

func (a *paymentClientAdapter) CreatePayment(ctx context.Context, bookingID string, amount float64) error {
	_, err := a.client.CreatePayment(ctx, &httpclient.PaymentRequest{
		BookingID: bookingID,
		Amount:    amount,
		Currency:  "RUB",
	})
	return err
}
