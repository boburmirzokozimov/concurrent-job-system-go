package worker

import (
	"fmt"
	"time"
)

type ExcelJob struct {
	BaseJob
	FilePath string
}

func NewExcelJob(filePath string, maxRetryCount int, prior Priority) *ExcelJob {
	return &ExcelJob{
		BaseJob: BaseJob{
			ID:            nextJobID(),
			MaxRetryCount: maxRetryCount,
			Priority:      prior,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		FilePath: filePath,
	}
}

func (j *ExcelJob) Process() error {
	LogJobStart(j)
	if j.ID%4 == 0 {
		LogJobFail(j)
		return fmt.Errorf("simulated failure on job %d", j.ID)
	}
	time.Sleep(2 * time.Second) // Simulate work
	LogJobSuccess(j)
	return nil
}
func (j *ExcelJob) Type() string {
	return "ExcelJob"
}
