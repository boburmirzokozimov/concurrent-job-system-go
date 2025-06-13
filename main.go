package main

import (
	config "concurrent-job-system/internal"
	"concurrent-job-system/models"
	"concurrent-job-system/worker"
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	_ "time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// Load config
	cfg := config.NewConfig()

	// Start worker pool
	storage, err := config.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	pool := worker.NewPool(5, storage) // 5 workers
	go pool.Start(ctx)

	simulate(ctx, sig, pool)

	<-sig
	log.Println("Shutdown signal received")
	cancel()
	pool.Wait()
	log.Println("Gracefully stopped")
}

func simulate(ctx context.Context, sig chan os.Signal, pool *worker.Pool) {
	// Simulate adding jobs
	for i := 0; i < 30; i++ {
		select {
		case <-ctx.Done():
			log.Printf("Stopped submission at job %d", i)
			break
		case <-sig:
			log.Printf("Stopped submission at job %d", i)
			break
		default:
			var job models.Processable
			if i%2 == 0 {
				job = worker.NewSimpleJob(3, models.Low)
			} else {
				job = worker.NewExcelJob("excel.path", 3, models.High)
			}
			pool.Submit(job, ctx)
		}
	}
}
