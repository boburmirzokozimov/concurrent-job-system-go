package main

import (
	"concurrent-job-system/internal/app"
	"concurrent-job-system/internal/container"
	"context"
	"log"
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
}
