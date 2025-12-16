package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"hotel-booking-system/internal/payment/domain"
	"hotel-booking-system/pkg/logger"

	"github.com/google/uuid"
)

type PaymentService struct {
	webhookURL string
}

func NewPaymentService(webhookURL string) *PaymentService {
	return &PaymentService{
		webhookURL: webhookURL,
	}
}

func (ps *PaymentService) ProcessPayment(ctx context.Context, req *domain.PaymentRequest) (*domain.PaymentResponse, error) {
	paymentID := uuid.New().String()

	response := &domain.PaymentResponse{
		PaymentID: paymentID,
		Status:    "processing",
		Message:   "payment is being processed",
	}

	go ps.processPaymentAsync(ctx, paymentID, req)

	return response, nil
}

func (ps *PaymentService) processPaymentAsync(ctx context.Context, paymentID string, req *domain.PaymentRequest) {
	time.Sleep(2 * time.Second)

	status := "paid"
	if req.Amount <= 0 {
		status = "failed"
	}

	webhook := domain.PaymentWebhook{
		PaymentID:   paymentID,
		BookingID:   req.BookingID,
		Status:      status,
		Amount:      req.Amount,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	if err := ps.sendWebhook(ctx, webhook); err != nil {
		logger.GetLogger().WithError(err).Error("failed to send payment webhook")
	}
}

func (ps *PaymentService) sendWebhook(ctx context.Context, webhook domain.PaymentWebhook) error {
	data, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ps.webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	logger.GetLogger().WithFields(map[string]interface{}{
		"payment_id": webhook.PaymentID,
		"booking_id": webhook.BookingID,
		"status":     webhook.Status,
	}).Info("payment webhook sent successfully")

	return nil
}
