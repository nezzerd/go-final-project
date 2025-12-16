package http

import (
	"encoding/json"
	"net/http"
	"time"

	"hotel-booking-system/internal/booking/domain"
	"hotel-booking-system/pkg/logger"
	"hotel-booking-system/pkg/metrics"

	"github.com/go-chi/chi/v5"
)

type BookingHandler struct {
	useCase domain.BookingUseCase
}

func NewBookingHandler(useCase domain.BookingUseCase) *BookingHandler {
	return &BookingHandler{useCase: useCase}
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/bookings").Observe(time.Since(start).Seconds())
	}()

	var booking domain.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		logger.GetLogger().WithError(err).Error("failed to decode request")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/bookings", "400").Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.useCase.CreateBooking(r.Context(), &booking); err != nil {
		logger.GetLogger().WithError(err).Error("failed to create booking")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/bookings", "500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/bookings", "201").Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/bookings/{id}").Observe(time.Since(start).Seconds())
	}()

	id := chi.URLParam(r, "id")
	booking, err := h.useCase.GetBooking(r.Context(), id)
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to get booking")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/bookings/{id}", "404").Inc()
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/bookings/{id}", "200").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) GetBookingsByUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/bookings/user/{userId}").Observe(time.Since(start).Seconds())
	}()

	userID := chi.URLParam(r, "userId")
	bookings, err := h.useCase.GetBookingsByUser(r.Context(), userID)
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to get bookings by user")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/bookings/user/{userId}", "500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/bookings/user/{userId}", "200").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) GetBookingsByHotel(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/bookings/hotel/{hotelId}").Observe(time.Since(start).Seconds())
	}()

	hotelID := chi.URLParam(r, "hotelId")
	bookings, err := h.useCase.GetBookingsByHotel(r.Context(), hotelID)
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to get bookings by hotel")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/bookings/hotel/{hotelId}", "500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/bookings/hotel/{hotelId}", "200").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

type PaymentWebhookRequest struct {
	PaymentID string  `json:"payment_id"`
	BookingID string  `json:"booking_id"`
	Status    string  `json:"status"`
	Amount    float64 `json:"amount,omitempty"`
}

func (h *BookingHandler) PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/webhooks/payment").Observe(time.Since(start).Seconds())
	}()

	var req PaymentWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLogger().WithError(err).Error("failed to decode webhook request")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/webhooks/payment", "400").Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.useCase.UpdatePaymentStatus(r.Context(), req.BookingID, req.Status); err != nil {
		logger.GetLogger().WithError(err).Error("failed to update payment status")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/webhooks/payment", "500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/webhooks/payment", "200").Inc()
	w.WriteHeader(http.StatusOK)
}
