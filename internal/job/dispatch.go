package job

import (
	"concurrent-job-system/internal/job/priority"
	"concurrent-job-system/internal/worker"
	"encoding/json"
	"reflect"
	"time"
)

type HandlerWithPayload interface {
	SetPayload(any)
	Process() error
}

type WrapperJob struct {
	BaseJob
	handler HandlerWithPayload
}

func (j *WrapperJob) Process() error {
	return j.handler.Process()
}

func (j *WrapperJob) PayloadOnly() interface{} {
	return j.BaseJob.Payload
}

func (j *WrapperJob) SetId(id int) error {
	j.ID = id
	return nil
}

func (j *WrapperJob) GetId() int {
	return j.ID
}

func (j *WrapperJob) GetStatus() string {
	return j.Status
}

func (j *WrapperJob) Type() string {
	return j.JobType
}

func (j *WrapperJob) GetPriority() priority.Priority {
	return priority.FromInt(j.Priority)
}

func (j *WrapperJob) GetRetries() int       { return j.Retries }
func (j *WrapperJob) GetMaxRetryCount() int { return j.MaxRetryCount }
func (j *WrapperJob) IncRetry()             { j.Retries++ }
func (j *WrapperJob) Base() *BaseJob        { return &j.BaseJob }

var pool *worker.Pool

func Dispatch(handler HandlerWithPayload, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	handler.SetPayload(payload)

	base := BaseJob{
		JobType:       reflect.TypeOf(handler).Elem().Name(),
		Payload:       string(data),
		Priority:      1,
		MaxRetryCount: 3,
		Status:        "queued",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	wrapper := &WrapperJob{
		BaseJob: base,
		handler: handler,
	}

	pool.Submit(wrapper)
	return nil
}
