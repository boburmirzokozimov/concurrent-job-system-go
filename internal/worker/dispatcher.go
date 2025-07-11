package worker

import (
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/logger"
	"concurrent-job-system/internal/metrics"
	"context"
	"sync"
	"time"
)

type Dispatcher struct {
	queues      *JobQueueSet
	executor    *JobExecutor
	stats       *JobStats
	logger      logger.ILogger
	busyWorkers sync.Map
}

func NewDispatcher(queues *JobQueueSet, executor *JobExecutor, log logger.ILogger) *Dispatcher {
	return &Dispatcher{
		queues:   queues,
		executor: executor,
		logger:   log,
	}
}

func (d *Dispatcher) Run(ctx context.Context, workerID int) {
	d.logger.Info("Worker %d started", workerID)
	for {
		select {
		case <-ctx.Done():
			d.logger.Info("Worker %d shutting down", workerID)
			return
		default:
			dequeue := d.queues.Dequeue()
			if dequeue == nil {
				dequeue = d.executor.LoadPending()
				if dequeue == nil {
					time.Sleep(10 * time.Second)
					continue
				}
			}
			d.busyWorkers.Store(workerID, true)
			d.logger.Debug("Worker %d picked job %d", workerID, dequeue.GetId())
			metrics.QueuedGauge.Dec()
			d.executor.Execute(dequeue, ctx)
			d.busyWorkers.Delete(workerID)
		}
	}
}
func (d *Dispatcher) Save(j job.IProcessable) (int, error) {
	return d.executor.Save(j)
}

func (d *Dispatcher) BusyWorkers() int {
	count := 0
	d.busyWorkers.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}
