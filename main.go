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
	storage := worker.NewFileJobStorage("jobs.json")
	pool := worker.NewPool(5, storage) // 5 workers
	go pool.Start(ctx)

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
			var job worker.Processable
			if i%2 == 0 {
				job = worker.NewSimpleJob(3, worker.Low)
			} else {
				job = worker.NewExcelJob("excel.path", 3, worker.High)
			}
			pool.Submit(job, ctx)
		}
	}

	<-sig
	log.Println("Shutdown signal received")
	cancel()
	pool.Wait()
	log.Println("Gracefully stopped")
}
