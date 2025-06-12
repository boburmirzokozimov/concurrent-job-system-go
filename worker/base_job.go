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
}

var globalJobID uint64

func nextJobID() int {
	return int(atomic.AddUint64(&globalJobID, 1))
}

type BaseJob struct {
	ID            int
	Retries       int
	MaxRetryCount int
	Priority      Priority
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
