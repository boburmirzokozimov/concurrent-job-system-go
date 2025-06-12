package worker

import (
	"sync/atomic"
	"time"
)

type Processable interface {
	GetRetries() int
	GetId() int
	GetMaxRetryCount() int
	Process() error
	IncRetry()
	Type() string
	GetPriority() Priority
	Base() *BaseJob
}

var globalJobID uint64

func nextJobID() int {
	return int(atomic.AddUint64(&globalJobID, 1))
}

type BaseJob struct {
	ID            int       `json:"id"`
	Retries       int       `json:"retries"`
	MaxRetryCount int       `json:"max_retry_count"`
	Priority      Priority  `json:"priority"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Status        string    `json:"status"` // pending, success, failed, canceled
}

func (j *BaseJob) GetRetries() int {
	return j.Retries
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

func (j *BaseJob) GetPriority() Priority {
	return j.Priority
}
func (j *BaseJob) Base() *BaseJob {
	return j
}
