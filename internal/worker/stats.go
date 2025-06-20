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

func (s *JobStats) Reset() {
	atomic.StoreUint64(&s.Total, 0)
	atomic.StoreUint64(&s.Success, 0)
	atomic.StoreUint64(&s.Failed, 0)
	s.Status = sync.Map{} // Clear map
}
func (s *JobStats) GetStatusCounts() map[string]int {
	counts := make(map[string]int)
	s.Status.Range(func(_, v any) bool {
		status := v.(string)
		counts[status]++
		return true
	})
	return counts
}
func (s *JobStats) PrintVerbose() {
	s.Print()
	fmt.Println("== Status Breakdown ==")
	for status, count := range s.GetStatusCounts() {
		fmt.Printf("  %s: %d\n", status, count)
	}
	fmt.Println()
}
