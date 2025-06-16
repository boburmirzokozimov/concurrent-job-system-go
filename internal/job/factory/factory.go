package factory

import (
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/job/types"
	"encoding/json"
	"fmt"
)

var payloadRegistry = map[string]func() job.Payload{
	"Simple": func() job.Payload { return &types.SimplePayloadJob{} },
	"Excel":  func() job.Payload { return &types.ExcelJobPayload{} },
}

func DeserializeJob(row job.BaseJob) (job.IProcessable, error) {
	creator, ok := payloadRegistry[row.JobType]
	if !ok {
		return nil, fmt.Errorf("unknown job type: %s", row.JobType)
	}

	payload := creator()
	if err := json.Unmarshal([]byte(row.Payload), payload); err != nil {
		return nil, err
	}

	return payload.ToProcessable(row), nil
}
