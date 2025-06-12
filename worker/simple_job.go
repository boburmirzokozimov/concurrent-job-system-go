package worker

import (
	"time"
)

type SimpleJob struct {
	BaseJob
}

func NewSimpleJob(maxRetryCount int, prior Priority) *SimpleJob {
	return &SimpleJob{
		BaseJob: BaseJob{
			ID:            nextJobID(),
			MaxRetryCount: maxRetryCount,
			Priority:      prior,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}
}

func (j *SimpleJob) Process() error {
	LogJobStart(j)
	time.Sleep(2 * time.Second) // Simulate work
	LogJobSuccess(j)
	return nil
}
func (j *SimpleJob) Type() string {
	return "SimpleJob"
}
