package types

import (
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/job/priority"
	"time"
)

type SimpleJob struct {
	job.BaseJob
}

func NewSimpleJob() *SimpleJob {
	return &SimpleJob{
		BaseJob: job.BaseJob{
			MaxRetryCount: 3,
			Priority:      int(priority.Low),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			JobType:       "Simple",
		},
	}
}

func (j *SimpleJob) Process() error {
	time.Sleep(2 * time.Second) // Simulate work
	return nil
}

func (j *SimpleJob) ToProcessable(_ job.BaseJob) job.IProcessable {
	return &SimpleJob{}
}

type SimplePayloadJob struct{}

func (j *SimpleJob) PayloadOnly() interface{} {
	return SimplePayloadJob{}
}
