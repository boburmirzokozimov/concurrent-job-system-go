package worker

import (
	"log"
	"time"
)

type Job struct {
	ID            int
	Retries       int
	MaxRetryCount int
}

func NewJob(id int, maxRetryCount int) Job {
	return Job{ID: id, MaxRetryCount: maxRetryCount}
}

func (j Job) process() error {
	log.Printf("Starting job %d", j.ID)
	time.Sleep(2 * time.Second) // Simulate work
	log.Printf("Finished job %d", j.ID)
	return nil
}
