package worker

import (
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/logger"
	"context"
	"time"
)

type Dispatcher struct {
	queues   *JobQueueSet
	executor *JobExecutor
	stats    *JobStats
	logger   logger.ILogger
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
					time.Sleep(100 * time.Millisecond)
					continue
				}
			}
			d.logger.Debug("Worker %d picked job %d", workerID, dequeue.GetId())
			d.executor.Execute(dequeue, ctx)
		}
	}
}
func (d *Dispatcher) Save(j job.IProcessable) (int, error) {
	return d.executor.Save(j)
}
