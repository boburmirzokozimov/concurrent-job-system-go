package worker

import (
	"concurrent-job-system/models"
	"time"
)

type SimpleJob struct {
	models.BaseJob
}

func NewSimpleJob(maxRetryCount int, prior models.Priority) *SimpleJob {
	return &SimpleJob{
		BaseJob: models.BaseJob{
			MaxRetryCount: maxRetryCount,
			Priority:      int(prior),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			JobType:       "Simple",
		},
	}
}

func (j *SimpleJob) Process() error {
	LogJobStart(j)
	time.Sleep(2 * time.Second) // Simulate work
	LogJobSuccess(j)
	return nil
}

func (j *SimpleJob) ToProcessable(_ models.BaseJob) models.Processable {
	return &SimpleJob{}
}

type SimplePayloadJob struct{}

func (j *SimpleJob) PayloadOnly() interface{} {
	return SimplePayloadJob{}
}
