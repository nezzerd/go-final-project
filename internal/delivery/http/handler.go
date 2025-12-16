package http

import (
	"encoding/json"
	"net/http"
	"time"

	"hotel-booking-system/internal/delivery/domain"
	"hotel-booking-system/internal/delivery/service"
	"hotel-booking-system/pkg/logger"
	"hotel-booking-system/pkg/metrics"
)

type DeliveryHandler struct {
	notifier service.Notifier
}

func NewDeliveryHandler(notifier service.Notifier) *DeliveryHandler {
	return &DeliveryHandler{
		notifier: notifier,
	}
}

func (h *DeliveryHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/notifications/send").Observe(time.Since(start).Seconds())
	}()

	var req domain.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLogger().WithError(err).Error("failed to decode request")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/notifications/send", "400").Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.notifier.SendNotification(&req); err != nil {
		logger.GetLogger().WithError(err).Error("failed to send notification")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/notifications/send", "500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/notifications/send", "200").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.SendNotificationResponse{
		Success: true,
		Message: "notification sent successfully",
	})
}
