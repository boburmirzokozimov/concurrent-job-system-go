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

	fmt.Printf(
		"\n\033[1m== Job Stats ==\033[0m\n"+
			"\033[36mTotal:\033[0m   %d\n"+
			"\033[32mSuccess:\033[0m %d\n"+
			"\033[31mFailed:\033[0m  %d\n\n",
		total, success, failed,
	)
}
