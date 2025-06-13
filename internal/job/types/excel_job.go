package types

import (
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/job/priority"
	"fmt"
	"time"
)

type ExcelJob struct {
	job.BaseJob
	FilePath string
}

func NewExcelJob(filePath string) *ExcelJob {
	return &ExcelJob{
		BaseJob: job.BaseJob{
			MaxRetryCount: 3,
			Priority:      int(priority.High),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			JobType:       "Excel",
		},
		FilePath: filePath,
	}
}

func (j *ExcelJob) Process() error {
	if j.ID%4 == 0 {
		return fmt.Errorf("simulated failure on job %d", j.ID)
	}
	time.Sleep(2 * time.Second) // Simulate work
	return nil
}

type ExcelJobPayload struct {
	FilePath string `json:"filePath"`
}

func (j *ExcelJobPayload) ToProcessable(base job.BaseJob) job.IProcessable {
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
