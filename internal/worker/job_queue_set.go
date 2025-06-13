package worker

import (
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/job/priority"
)

type JobQueueSet struct {
	high   chan job.IProcessable
	normal chan job.IProcessable
	low    chan job.IProcessable
}

func NewJobQueueSet(bufferSize int) *JobQueueSet {
	return &JobQueueSet{
		high:   make(chan job.IProcessable, bufferSize),
		normal: make(chan job.IProcessable, bufferSize),
		low:    make(chan job.IProcessable, bufferSize),
	}
}
func (q *JobQueueSet) Enqueue(j job.IProcessable) {
	switch j.GetPriority() {
	case priority.High:
		q.high <- j
	case priority.Normal:
		q.normal <- j
	case priority.Low:
		q.low <- j
	}
}

func (q *JobQueueSet) Dequeue() job.IProcessable {
	select {
	case j := <-q.high:
		return j
	default:
	}

	select {
	case j := <-q.normal:
		return j
	default:
	}

	select {
	case j := <-q.low:
		return j
	default:
	}

	return nil
}
