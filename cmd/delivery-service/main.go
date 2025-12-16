package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpHandler "hotel-booking-system/internal/delivery/http"
	"hotel-booking-system/internal/delivery/service"
	"hotel-booking-system/pkg/logger"
	"hotel-booking-system/pkg/tracing"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()

	logger.Init(os.Getenv("LOG_LEVEL"))
	log := logger.GetLogger()

	tp, err := tracing.InitTracer("delivery-service", os.Getenv("JAEGER_ENDPOINT"))
	if err != nil {
		log.WithError(err).Fatal("failed to init tracer")
	}
	defer tracing.Shutdown(context.Background(), tp)

	deliveryService, err := service.NewDeliveryService(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.WithError(err).Fatal("failed to create delivery service")
	}

	handler := httpHandler.NewDeliveryHandler(deliveryService)
	router := httpHandler.SetupRoutes(handler)

	httpPort := os.Getenv("DELIVERY_SERVICE_PORT")
	if httpPort == "" {
		httpPort = "8084"
	}

	go func() {
		prometheusPort := os.Getenv("PROMETHEUS_PORT")
		http.Handle("/metrics", promhttp.Handler())
		log.Infof("starting prometheus metrics on port %s", prometheusPort)
		http.ListenAndServe(":"+prometheusPort, nil)
	}()

	server := &http.Server{
		Addr:    ":" + httpPort,
		Handler: router,
	}

	go func() {
		log.Infof("starting HTTP server on port %s", httpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("failed to start HTTP server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down delivery service")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("failed to shutdown server gracefully")
	}
}
