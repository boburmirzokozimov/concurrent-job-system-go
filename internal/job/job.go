package job

import (
	"concurrent-job-system/internal/job/priority"
	"time"
)

type IProcessable interface {
	GetRetries() int
	GetId() int
	GetMaxRetryCount() int
	Process() error
	IncRetry()
	Type() string
	GetPriority() priority.Priority
	Base() *BaseJob
	PayloadOnly() interface{}
	SetId(id int) error
	GetStatus() string
}

type Payload interface {
	ToProcessable(base BaseJob) IProcessable
}

type BaseJob struct {
	ID            int `gorm:"primaryKey;autoIncrement"`
	JobType       string
	Payload       string
	Priority      int
	MaxRetryCount int
	Retries       int
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (j *BaseJob) GetRetries() int {
	return j.Retries
}

func (j *BaseJob) GetStatus() string {
	return j.Status
}

func (j *BaseJob) GetMaxRetryCount() int {
	return j.MaxRetryCount
}

func (j *BaseJob) GetId() int {
	return j.ID
}

func (j *BaseJob) IncRetry() {
	j.Retries++
	j.UpdatedAt = time.Now()
}
func (j *BaseJob) SetId(id int) error {
	j.ID = id
	return nil
}

func (j *BaseJob) GetPriority() priority.Priority {
	return priority.FromInt(j.Priority)
}
func (j *BaseJob) Base() *BaseJob {
	return j
}
func (j *BaseJob) Type() string {
	return j.JobType
}
