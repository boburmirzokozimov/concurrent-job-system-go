package app

import (
	"concurrent-job-system/cmd/api"
	"concurrent-job-system/internal/container"
	"context"
	"log"
	"net/http"
)

type App struct {
	container *container.Container
}

func NewApp(container *container.Container) *App {
	return &App{
		container: container,
	}
}

func (a *App) Run(ctx context.Context) error {
	//Start a pool
	go a.container.Pool.Start(ctx)

	//Simulate jobs
	go func(con *container.Container) {
		handler := api.NewHandler(con)
		log.Println("📡 Starting REST API server at :8080")
		if err := http.ListenAndServe(":8080", api.NewRouter(handler)); err != nil {
			log.Fatalf("API server failed: %v", err)
		}
	}(a.container)
	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Wait for workers to finish
	a.container.Pool.Wait()
	log.Println("Gracefully stopped")

	return nil
}
