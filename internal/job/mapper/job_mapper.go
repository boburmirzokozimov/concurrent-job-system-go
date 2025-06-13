package mapper

import (
	"concurrent-job-system/internal/job"
	"encoding/json"
)

func ToBaseJob(j job.IProcessable) job.BaseJob {
	payload, _ := json.Marshal(j.PayloadOnly())
	return job.BaseJob{
		JobType:       j.Type(),
		Payload:       string(payload),
		Priority:      int(j.GetPriority()),
		MaxRetryCount: j.GetMaxRetryCount(),
		Retries:       j.GetRetries(),
		Status:        "pending",
		CreatedAt:     j.Base().CreatedAt,
		UpdatedAt:     j.Base().UpdatedAt,
	}
}
