package worker

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type JobStats struct {
	Total   uint64
	Success uint64
	Failed  uint64

	Status sync.Map
}

func (s *JobStats) IncSuccess() {
	atomic.AddUint64(&s.Success, 1)
}
func (s *JobStats) IncFailed() {
	atomic.AddUint64(&s.Failed, 1)
}

func (s *JobStats) IncTotal() {
	atomic.AddUint64(&s.Total, 1)
}

func (s *JobStats) RecordStatus(jobId int, status string) {
	s.Status.Store(jobId, status)
}

func (s *JobStats) Print() {
	total := atomic.LoadUint64(&s.Total)
	success := atomic.LoadUint64(&s.Success)
	failed := atomic.LoadUint64(&s.Failed)

	fmt.Printf("== Job Stats ==\nTotal: %d\nSuccess: %d\nFailed: %d\n", total, success, failed)
}
