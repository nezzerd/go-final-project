package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"hotel-booking-system/internal/payment/domain"
	"hotel-booking-system/pkg/logger"
	"hotel-booking-system/pkg/metrics"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, req *domain.PaymentRequest) (*domain.PaymentResponse, error)
}

type PaymentHandler struct {
	paymentService PaymentService
}

func NewPaymentHandler(paymentService PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/payments").Observe(time.Since(start).Seconds())
	}()

	var req domain.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLogger().WithError(err).Error("failed to decode request")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/payments", "400").Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Currency == "" {
		req.Currency = "RUB"
	}

	response, err := h.paymentService.ProcessPayment(r.Context(), &req)
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to process payment")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/payments", "500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/payments", "202").Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}
