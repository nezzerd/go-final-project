package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpHandler "hotel-booking-system/internal/payment/delivery/http"
	"hotel-booking-system/internal/payment/service"
	"hotel-booking-system/pkg/logger"
	"hotel-booking-system/pkg/tracing"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()

	logger.Init(os.Getenv("LOG_LEVEL"))
	log := logger.GetLogger()

	tp, err := tracing.InitTracer("payment-service", os.Getenv("JAEGER_ENDPOINT"))
	if err != nil {
		log.WithError(err).Fatal("failed to init tracer")
	}
	defer tracing.Shutdown(context.Background(), tp)

	webhookURL := os.Getenv("BOOKING_WEBHOOK_URL")
	if webhookURL == "" {
		webhookURL = "http://booking-service:8082/api/webhooks/payment"
	}

	paymentService := service.NewPaymentService(webhookURL)
	handler := httpHandler.NewPaymentHandler(paymentService)
	router := httpHandler.SetupRoutes(handler)

	httpPort := os.Getenv("PAYMENT_SERVICE_PORT")
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

	log.Info("shutting down payment service")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("failed to shutdown server gracefully")
	}
}
