package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	httpHandler "hotel-booking-system/internal/booking/delivery/http"
	"hotel-booking-system/internal/booking/repository"
	"hotel-booking-system/internal/booking/usecase"
	"hotel-booking-system/pkg/database"
	"hotel-booking-system/pkg/hotelclient"
	"hotel-booking-system/pkg/httpclient"
	"hotel-booking-system/pkg/kafka"
	"hotel-booking-system/pkg/logger"
	"hotel-booking-system/pkg/tracing"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()

	logger.Init(os.Getenv("LOG_LEVEL"))
	log := logger.GetLogger()

	tp, err := tracing.InitTracer("booking-service", os.Getenv("JAEGER_ENDPOINT"))
	if err != nil {
		log.WithError(err).Fatal("failed to init tracer")
	}
	defer tracing.Shutdown(context.Background(), tp)

	dbCfg := database.Config{
		Host:     os.Getenv("BOOKING_DB_HOST"),
		Port:     os.Getenv("BOOKING_DB_PORT"),
		User:     os.Getenv("BOOKING_DB_USER"),
		Password: os.Getenv("BOOKING_DB_PASSWORD"),
		DBName:   os.Getenv("BOOKING_DB_NAME"),
	}

	db, err := database.NewPostgresConnection(dbCfg)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}
	defer db.Close()

	bookingRepo := repository.NewPostgresBookingRepository(db)

	hotelClient, err := hotelclient.NewHotelClient(os.Getenv("HOTEL_SERVICE_HOST"))
	if err != nil {
		log.WithError(err).Fatal("failed to create hotel client")
	}
	defer hotelClient.Close()

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	producer := kafka.NewProducer(brokers, os.Getenv("KAFKA_TOPIC_BOOKING_CREATED"))
	defer producer.Close()

	var paymentClient usecase.PaymentClient
	paymentServiceURL := os.Getenv("PAYMENT_SERVICE_URL")
	if paymentServiceURL != "" {
		paymentHTTPClient := httpclient.NewPaymentClient(paymentServiceURL)
		var paymentClientInterface usecase.PaymentClientInterface = paymentHTTPClient
		paymentClient = usecase.NewPaymentClientAdapter(paymentClientInterface)
	}

	bookingUseCase := usecase.NewBookingUseCase(bookingRepo, hotelClient, producer, paymentClient)

	httpPort := os.Getenv("BOOKING_SERVICE_PORT")

	go func() {
		handler := httpHandler.NewBookingHandler(bookingUseCase)
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

	log.Info("shutting down booking service")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = ctx
}
