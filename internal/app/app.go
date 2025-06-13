package app

import (
	"concurrent-job-system/internal/container"
	"context"
	"log"
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
	go a.simulateJobs(ctx)

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Wait for workers to finish
	a.container.Pool.Wait()
	log.Println("Gracefully stopped")

	return nil
}

func (a *App) simulateJobs(ctx context.Context) {
	for i := 0; i < 30; i++ {
		select {
		case <-ctx.Done():
			log.Printf("Stopped job submission at job %d", i)
			return
		default:
			var job = a.container.MakeJob(i)
			a.container.Pool.Submit(job)
		}
	}
}
