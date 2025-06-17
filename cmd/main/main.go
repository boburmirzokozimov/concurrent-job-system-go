package main

import (
	"concurrent-job-system/internal/app"
	"concurrent-job-system/internal/container"
	"concurrent-job-system/internal/metrics"
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	c, err := container.NewContainer()
	if err != nil {
		log.Fatal("Failed to initialize app container:", err)
	}

	application := app.NewApp(c)
	err = application.Run(ctx)
	if err != nil {
		log.Fatal("Got some error:", err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	metrics.Init()
	StartMetricsServer()
}

func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Prometheus metrics exposed at :2112/metrics")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()
}
