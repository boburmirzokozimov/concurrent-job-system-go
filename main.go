package main

import (
	"concurrent-job-system/worker"
	"context"
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

	// Start worker pool
	pool := worker.NewPool(5) // 5 workers
	go pool.Start(ctx)

	// Simulate adding jobs
	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			log.Printf("Stopped submission at job %d", i)
			break
		case <-sig:
			log.Printf("Stopped submission at job %d", i)
			break
		default:
			job := worker.NewJob(i, 3)
			pool.Submit(job, ctx)
		}
	}

	<-sig
	log.Println("Shutdown signal received")
	cancel()
	pool.Wait()
	log.Println("Gracefully stopped")
}
