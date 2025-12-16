package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"hotel-booking-system/internal/booking/domain"
	"hotel-booking-system/internal/notification/service"
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

	tp, err := tracing.InitTracer("notification-service", os.Getenv("JAEGER_ENDPOINT"))
	if err != nil {
		log.WithError(err).Fatal("failed to init tracer")
	}
	defer tracing.Shutdown(context.Background(), tp)

	deliveryServiceURL := os.Getenv("DELIVERY_SERVICE_URL")
	if deliveryServiceURL == "" {
		deliveryServiceURL = "http://delivery-service:8084"
	}

	hotelServiceURL := os.Getenv("HOTEL_SERVICE_URL")
	if hotelServiceURL == "" {
		hotelServiceURL = "http://hotel-service:8081"
	}

	deliveryClient := httpclient.NewDeliveryClient(deliveryServiceURL)
	hotelClient := httpclient.NewHotelHTTPClient(hotelServiceURL)

	var deliveryClientInterface service.DeliveryClient = deliveryClient
	var hotelClientInterface service.HotelClient = hotelClient

	notificationService := service.NewNotificationService(deliveryClientInterface, hotelClientInterface)

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	consumer := kafka.NewConsumer(
		brokers,
		os.Getenv("KAFKA_TOPIC_BOOKING_CREATED"),
		os.Getenv("KAFKA_GROUP_ID"),
	)
	defer consumer.Close()

	go func() {
		prometheusPort := os.Getenv("PROMETHEUS_PORT")
		http.Handle("/metrics", promhttp.Handler())
		log.Infof("starting prometheus metrics on port %s", prometheusPort)
		http.ListenAndServe(":"+prometheusPort, nil)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Info("starting kafka consumer")
		consumer.ReadMessage(ctx, func(data []byte) error {
			var event domain.BookingEvent
			if err := kafka.UnmarshalMessage(data, &event); err != nil {
				log.WithError(err).Error("failed to unmarshal booking event")
				return err
			}

			log.WithField("booking_id", event.BookingID).Info("received booking event")

			if err := notificationService.ProcessBookingEvent(ctx, event); err != nil {
				log.WithError(err).Error("failed to process booking event")
			}

			return nil
		})
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down notification service")
	cancel()
	time.Sleep(2 * time.Second)
}
