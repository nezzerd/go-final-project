package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(handler *BookingHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api", func(r chi.Router) {
		r.Route("/bookings", func(r chi.Router) {
			r.Post("/", handler.CreateBooking)
			r.Get("/{id}", handler.GetBooking)
			r.Get("/user/{userId}", handler.GetBookingsByUser)
			r.Get("/hotel/{hotelId}", handler.GetBookingsByHotel)
		})

		r.Route("/webhooks", func(r chi.Router) {
			r.Post("/payment", handler.PaymentWebhook)
		})
	})

	return r
}
