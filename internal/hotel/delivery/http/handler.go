package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"hotel-booking-system/internal/hotel/domain"
	"hotel-booking-system/pkg/logger"
	"hotel-booking-system/pkg/metrics"

	"github.com/go-chi/chi/v5"
)

type HotelHandler struct {
	useCase domain.HotelUseCase
}

func NewHotelHandler(useCase domain.HotelUseCase) *HotelHandler {
	return &HotelHandler{useCase: useCase}
}

func (h *HotelHandler) CreateHotel(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/hotels").Observe(time.Since(start).Seconds())
	}()

	var hotel domain.Hotel
	if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
		logger.GetLogger().WithError(err).Error("failed to decode request")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels", "400").Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.useCase.CreateHotel(r.Context(), &hotel); err != nil {
		logger.GetLogger().WithError(err).Error("failed to create hotel")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels", "500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels", "201").Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(hotel)
}

func (h *HotelHandler) GetHotel(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/hotels/{id}").Observe(time.Since(start).Seconds())
	}()

	id := chi.URLParam(r, "id")
	hotel, err := h.useCase.GetHotel(r.Context(), id)
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to get hotel")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels/{id}", "404").Inc()
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels/{id}", "200").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotel)
}

func (h *HotelHandler) GetHotels(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/hotels").Observe(time.Since(start).Seconds())
	}()

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	hotels, err := h.useCase.GetHotels(r.Context(), limit, offset)
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to get hotels")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels", "500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels", "200").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotels)
}

func (h *HotelHandler) GetHotelWithRooms(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/hotels/{id}/rooms").Observe(time.Since(start).Seconds())
	}()

	id := chi.URLParam(r, "id")
	hotelWithRooms, err := h.useCase.GetHotelWithRooms(r.Context(), id)
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to get hotel with rooms")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels/{id}/rooms", "404").Inc()
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels/{id}/rooms", "200").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotelWithRooms)
}

func (h *HotelHandler) UpdateHotel(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/hotels/{id}").Observe(time.Since(start).Seconds())
	}()

	id := chi.URLParam(r, "id")
	var hotel domain.Hotel
	if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
		logger.GetLogger().WithError(err).Error("failed to decode request")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels/{id}", "400").Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hotel.ID = id

	if err := h.useCase.UpdateHotel(r.Context(), &hotel); err != nil {
		logger.GetLogger().WithError(err).Error("failed to update hotel")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels/{id}", "500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/hotels/{id}", "200").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotel)
}

func (h *HotelHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, "/api/rooms").Observe(time.Since(start).Seconds())
	}()

	var room domain.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		logger.GetLogger().WithError(err).Error("failed to decode request")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/rooms", "400").Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.useCase.CreateRoom(r.Context(), &room); err != nil {
		logger.GetLogger().WithError(err).Error("failed to create room")
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/rooms", "500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, "/api/rooms", "201").Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}
