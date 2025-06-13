package worker

import (
	"concurrent-job-system/models"
	"fmt"
	"time"
)

type ExcelJob struct {
	models.BaseJob
	FilePath string
}

func NewExcelJob(filePath string, maxRetryCount int, prior models.Priority) *ExcelJob {
	return &ExcelJob{
		BaseJob: models.BaseJob{
			MaxRetryCount: maxRetryCount,
			Priority:      int(prior),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			JobType:       "Excel",
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

type ExcelJobPayload struct {
	FilePath string `json:"filePath"`
}

func (j *ExcelJobPayload) ToProcessable(base models.BaseJob) models.Processable {
	return &ExcelJob{
		BaseJob:  base,
		FilePath: j.FilePath,
	}
}

func (j *ExcelJob) PayloadOnly() interface{} {
	return ExcelJobPayload{
		FilePath: j.FilePath,
	}
}
