package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(handler *HotelHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api", func(r chi.Router) {
		r.Route("/hotels", func(r chi.Router) {
			r.Get("/", handler.GetHotels)
			r.Post("/", handler.CreateHotel)
			r.Get("/{id}", handler.GetHotel)
			r.Put("/{id}", handler.UpdateHotel)
			r.Get("/{id}/rooms", handler.GetHotelWithRooms)
		})

		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", handler.CreateRoom)
		})
	})

	return r
}
