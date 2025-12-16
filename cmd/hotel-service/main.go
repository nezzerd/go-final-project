package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpHandler "hotel-booking-system/internal/hotel/delivery/http"
	"hotel-booking-system/internal/hotel/repository"
	"hotel-booking-system/internal/hotel/usecase"
	"hotel-booking-system/pkg/database"
	"hotel-booking-system/pkg/logger"
	"hotel-booking-system/pkg/tracing"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()

	logger.Init(os.Getenv("LOG_LEVEL"))
	log := logger.GetLogger()

	tp, err := tracing.InitTracer("hotel-service", os.Getenv("JAEGER_ENDPOINT"))
	if err != nil {
		log.WithError(err).Fatal("failed to init tracer")
	}
	defer tracing.Shutdown(context.Background(), tp)

	dbCfg := database.Config{
		Host:     os.Getenv("HOTEL_DB_HOST"),
		Port:     os.Getenv("HOTEL_DB_PORT"),
		User:     os.Getenv("HOTEL_DB_USER"),
		Password: os.Getenv("HOTEL_DB_PASSWORD"),
		DBName:   os.Getenv("HOTEL_DB_NAME"),
	}

	db, err := database.NewPostgresConnection(dbCfg)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}
	defer db.Close()

	hotelRepo := repository.NewPostgresHotelRepository(db)
	roomRepo := repository.NewPostgresRoomRepository(db)
	hotelUseCase := usecase.NewHotelUseCase(hotelRepo, roomRepo)

	httpPort := os.Getenv("HOTEL_SERVICE_PORT")

	go func() {
		handler := httpHandler.NewHotelHandler(hotelUseCase)
		router := httpHandler.SetupRoutes(handler)

		log.Infof("starting HTTP server on port %s", httpPort)
		if err := http.ListenAndServe(":"+httpPort, router); err != nil {
			log.WithError(err).Fatal("failed to start HTTP server")
		}
	}()

	go func() {
		prometheusPort := os.Getenv("PROMETHEUS_PORT")
		http.Handle("/metrics", promhttp.Handler())
		log.Infof("starting prometheus metrics on port %s", prometheusPort)
		http.ListenAndServe(":"+prometheusPort, nil)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down hotel service")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = ctx
}
